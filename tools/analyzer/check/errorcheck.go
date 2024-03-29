package check

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"golang.org/x/tools/go/loader"
)

var errorType *types.Interface

func init() {
	errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

}

var (
	// ErrNoGoFiles is returned when CheckPackage is run on a package with no Go source files
	ErrNoGoFiles = errors.New("package contains no go source files")
)

// UncheckedError indicates the position of an unchecked error return.
type UncheckedError struct {
	Pos      token.Position
	Line     string
	FuncName string
}

// UncheckedErrors is returned from the CheckPackage function if the package contains
// any unchecked errors.
// Errors should be appended using the Append method, which is safe to use concurrently.
type UncheckedErrors struct {
	mu sync.Mutex

	// Errors is a list of all the unchecked errors in the package.
	// Printing an error reports its position within the file and the contents of the line.
	Errors []UncheckedError
}

func (e *UncheckedErrors) Append(errors ...UncheckedError) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Errors = append(e.Errors, errors...)
}

func (e *UncheckedErrors) Error() string {
	return fmt.Sprintf("%d unchecked errors", len(e.Errors))
}

// Len is the number of elements in the collection.
func (e *UncheckedErrors) Len() int { return len(e.Errors) }

// Swap swaps the elements with indexes i and j.
func (e *UncheckedErrors) Swap(i, j int) { e.Errors[i], e.Errors[j] = e.Errors[j], e.Errors[i] }

type byName struct{ *UncheckedErrors }

// Less reports whether the element with index i should sort before the element with index j.
func (e byName) Less(i, j int) bool {
	ei, ej := e.Errors[i], e.Errors[j]

	pi, pj := ei.Pos, ej.Pos

	if pi.Filename != pj.Filename {
		return pi.Filename < pj.Filename
	}
	if pi.Line != pj.Line {
		return pi.Line < pj.Line
	}
	if pi.Column != pj.Column {
		return pi.Column < pj.Column
	}

	return ei.Line < ej.Line
}

type Checker struct {
	// ignore is a map of package names to regular expressions. Identifiers from a package are
	// checked against its regular expressions and if any of the expressions match the call
	// is not checked.
	Ignore map[string]*regexp.Regexp

	// If blank is true then assignments to the blank identifier are also considered to be
	// ignored errors.
	Blank bool

	// If asserts is true then ignored type assertion results are also checked
	Asserts bool

	// build tags
	Tags []string

	Verbose bool

	// If true, checking of of _test.go files is disabled
	WithoutTests bool

	exclude map[string]bool
}

func NewChecker() *Checker {
	c := Checker{}
	c.SetExclude(map[string]bool{})
	return &c
}

func (c *Checker) SetExclude(l map[string]bool) {
	// Default exclude for stdlib functions
	c.exclude = map[string]bool{
		"math/rand.Read":         true,
		"(*math/rand.Rand).Read": true,

		"(*bytes.Buffer).Write":       true,
		"(*bytes.Buffer).WriteByte":   true,
		"(*bytes.Buffer).WriteRune":   true,
		"(*bytes.Buffer).WriteString": true,
	}
	for k := range l {
		c.exclude[k] = true
	}
}

func (c *Checker) logf(msg string, args ...interface{}) {
	if c.Verbose {
		fmt.Fprintf(os.Stderr, msg+"\n", args...)
	}
}

func (c *Checker) load(paths ...string) (*loader.Program, error) {
	ctx := build.Default
	ctx.BuildTags = append(ctx.BuildTags, c.Tags...)
	loadcfg := loader.Config{
		Build: &ctx,
	}
	rest, err := loadcfg.FromArgs(paths, !c.WithoutTests)
	if err != nil {
		return nil, fmt.Errorf("could not parse arguments: %s", err)
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("unhandled extra arguments: %v", rest)
	}

	return loadcfg.Load()
}

func ErrorCheck(projectPath string) []string {
	errorcheck := NewChecker()
	errorcheck.Asserts = false
	errorcheck.Blank = false
	errorcheck.WithoutTests = true
	errorcheck.Verbose = true

	uncheckedErrors, err := errorcheck.CheckPackages(projectPath)
	if err != nil {
		return nil
	}

	return uncheckedErrors
}

// CheckPackages checks packages for errors.
func (c *Checker) CheckPackages(paths ...string) ([]string, error) {
	program, err := c.load(paths...)
	if err != nil {
		return nil, fmt.Errorf("could not type check: %s", err)
	}

	var wg sync.WaitGroup
	u := &UncheckedErrors{}
	for _, pkgInfo := range program.InitialPackages() {
		if pkgInfo.Pkg.Path() == "unsafe" { // not a real package
			continue
		}

		wg.Add(1)

		go func(pkgInfo *loader.PackageInfo) {
			defer wg.Done()
			c.logf("Checking %s", pkgInfo.Pkg.Path())

			v := &visitor{
				prog:    program,
				pkg:     pkgInfo,
				ignore:  c.Ignore,
				blank:   c.Blank,
				asserts: c.Asserts,
				lines:   make(map[string][]string),
				exclude: c.exclude,
				errors:  []UncheckedError{},
			}

			for _, astFile := range v.pkg.Files {
				ast.Walk(v, astFile)
			}
			u.Append(v.errors...)
		}(pkgInfo)
	}

	wg.Wait()
	if u.Len() > 0 {
		sort.Sort(byName{u})

		return reportUncheckedErrors(u, c.Verbose), nil
	}
	return nil, nil
}

// visitor implements the errcheck algorithm
type visitor struct {
	prog                             *loader.Program
	pkg                              *loader.PackageInfo
	ignore                           map[string]*regexp.Regexp
	blank                            bool
	asserts                          bool
	lines                            map[string][]string
	exclude                          map[string]bool
	parantBlockStart, parentBlockEnd token.Pos
	blockStart, blockEnd             token.Pos

	check       bool // 开始检查error标记
	check_error UncheckedError
	check_end   int

	errors []UncheckedError
}

func (v *visitor) fullName(call *ast.CallExpr) (string, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	fn, ok := v.pkg.ObjectOf(sel.Sel).(*types.Func)
	if !ok {
		// Shouldn't happen, but be paranoid
		return "", false
	}
	// The name is fully qualified by the import path, possible type,
	// function/method name and pointer receiver.
	//
	// TODO(dh): vendored packages will have /vendor/ in their name,
	// thus not matching vendored standard library packages. If we
	// want to support vendored stdlib packages, we need to implement
	// FullName with our own logic.
	return fn.FullName(), true
}

func (v *visitor) excludeCall(call *ast.CallExpr) bool {
	if name, ok := v.fullName(call); ok {
		return v.exclude[name]
	}

	return false
}

func (v *visitor) ignoreCall(call *ast.CallExpr) bool {
	if v.excludeCall(call) {
		return true
	}

	// Try to get an identifier.
	// Currently only supports simple expressions:
	//     1. f()
	//     2. x.y.f()
	var id *ast.Ident
	switch exp := call.Fun.(type) {
	case (*ast.Ident):
		id = exp
	case (*ast.SelectorExpr):
		id = exp.Sel
	default:
		// eg: *ast.SliceExpr, *ast.IndexExpr
	}

	if id == nil {
		return false
	}

	// If we got an identifier for the function, see if it is ignored
	if re, ok := v.ignore[""]; ok && re.MatchString(id.Name) {
		return true
	}

	if obj := v.pkg.Uses[id]; obj != nil {
		if pkg := obj.Pkg(); pkg != nil {
			if re, ok := v.ignore[pkg.Path()]; ok {
				return re.MatchString(id.Name)
			}

			// if current package being considered is vendored, check to see if it should be ignored based
			// on the unvendored path.
			if nonVendoredPkg, ok := nonVendoredPkgPath(pkg.Path()); ok {
				if re, ok := v.ignore[nonVendoredPkg]; ok {
					return re.MatchString(id.Name)
				}
			}
		}
	}

	return false
}

// nonVendoredPkgPath returns the unvendored version of the provided package path (or returns the provided path if it
// does not represent a vendored path). The second return value is true if the provided package was vendored, false
// otherwise.
func nonVendoredPkgPath(pkgPath string) (string, bool) {
	lastVendorIndex := strings.LastIndex(pkgPath, "/vendor/")
	if lastVendorIndex == -1 {
		return pkgPath, false
	}
	return pkgPath[lastVendorIndex+len("/vendor/"):], true
}

// errorsByArg returns a slice s such that
// len(s) == number of return types of call
// s[i] == true iff return type at position i from left is an error type
func (v *visitor) errorsByArg(call *ast.CallExpr) []bool {
	switch t := v.pkg.Types[call].Type.(type) {
	case *types.Named:
		// Single return
		return []bool{isErrorType(t)}
	case *types.Pointer:
		// Single return via pointer
		return []bool{isErrorType(t)}
	case *types.Tuple:
		// Multiple returns
		s := make([]bool, t.Len())
		for i := 0; i < t.Len(); i++ {
			switch et := t.At(i).Type().(type) {
			case *types.Named:
				// Single return
				s[i] = isErrorType(et)
			case *types.Pointer:
				// Single return via pointer
				s[i] = isErrorType(et)
			default:
				s[i] = false
			}
		}
		return s
	}
	return []bool{false}
}

func (v *visitor) callReturnsError(call *ast.CallExpr) bool {
	if v.isRecover(call) {
		return true
	}
	for _, isError := range v.errorsByArg(call) {
		if isError {
			return true
		}
	}
	return false
}

// isRecover returns true if the given CallExpr is a call to the built-in recover() function.
func (v *visitor) isRecover(call *ast.CallExpr) bool {
	if fun, ok := call.Fun.(*ast.Ident); ok {
		if _, ok := v.pkg.Uses[fun].(*types.Builtin); ok {
			return fun.Name == "recover"
		}
	}
	return false
}

func (v *visitor) addErrorAtPosition(position token.Pos, call *ast.CallExpr) {
	pos := v.prog.Fset.Position(position)
	lines, ok := v.lines[pos.Filename]
	if !ok {
		lines = readfile(pos.Filename)
		v.lines[pos.Filename] = lines
	}

	line := "??"
	if pos.Line-1 < len(lines) {
		line = strings.TrimSpace(lines[pos.Line-1])
	}

	var name string
	if call != nil {
		name, _ = v.fullName(call)
	}

	v.errors = append(v.errors, UncheckedError{pos, line, name})
}

func (v *visitor) walkBlock(start, end token.Pos) {
	// 上一层block error没有处理
	if v.blockStart != 0 && v.check {

	}
	v.blockStart, v.blockEnd = start, end
}

func (v *visitor) startErrorAtPosition(position token.Pos, name string) {
	pos := v.prog.Fset.Position(position)
	lines, ok := v.lines[pos.Filename]
	if !ok {
		lines = readfile(pos.Filename)
		v.lines[pos.Filename] = lines
	}

	line := "??"
	if pos.Line-1 < len(lines) {
		line = strings.TrimSpace(lines[pos.Line-1])
	}
	v.check = true // 开始检查
	v.check_error = UncheckedError{pos, line, name}
}

func (v *visitor) endErrorAtPosition(handle bool) {
	v.check = false // 结束检查
	if !handle {
		v.errors = append(v.errors, v.check_error)
	}
}

func readfile(filename string) []string {
	var f, err = os.Open(filename)
	if err != nil {
		return nil
	}

	var lines []string
	var scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch stmt := node.(type) {
	case *ast.ExprStmt:
		fmt.Println("expr")
		if call, ok := stmt.X.(*ast.CallExpr); ok {
			if !v.ignoreCall(call) && v.callReturnsError(call) {
				v.addErrorAtPosition(call.Lparen, call)
			}
		}
	case *ast.GoStmt:
		if !v.ignoreCall(stmt.Call) && v.callReturnsError(stmt.Call) {
			v.addErrorAtPosition(stmt.Call.Lparen, stmt.Call)
		}
	case *ast.DeferStmt:
		if !v.ignoreCall(stmt.Call) && v.callReturnsError(stmt.Call) {
			v.addErrorAtPosition(stmt.Call.Lparen, stmt.Call)
		}
	case *ast.ReturnStmt: // return err
		fmt.Println("returnStmt results: ", stmt.Results)
	case *ast.BlockStmt:
		v.walkBlock(stmt.Lbrace, stmt.Rbrace)
		fmt.Println("BlockStmt L: ", stmt.Lbrace, ", R:", stmt.Rbrace)
	case *ast.AssignStmt:
		fmt.Println("assign")
		if len(stmt.Rhs) == 1 {
			fmt.Println("czx@@@ assign1")
			// single value on rhs; check against lhs identifiers
			if call, ok := stmt.Rhs[0].(*ast.CallExpr); ok {
				fmt.Println("czx@@@ assign2")
				/*
					if !v.blank {
						break
					}
				*/
				fmt.Println("czx@@@ assign3")
				if v.ignoreCall(call) {
					break
				}
				fmt.Println("czx@@@ assign4")
				isError := v.errorsByArg(call)
				for i := 0; i < len(stmt.Lhs); i++ {
					fmt.Println("czx@@@ isError: ", isError[i])
					if id, ok := stmt.Lhs[i].(*ast.Ident); ok {
						// We shortcut calls to recover() because errorsByArg can't
						// check its return types for errors since it returns interface{}.
						if id.Name != "_" && isError[i] {
							v.startErrorAtPosition(id.NamePos, id.Name)
						}
					}
				}
			} else if assert, ok := stmt.Rhs[0].(*ast.TypeAssertExpr); ok {
				fmt.Println("czx@@@ assign3")
				if !v.asserts {
					break
				}
				if assert.Type == nil {
					// type switch
					break
				}
				if len(stmt.Lhs) < 2 {
					// assertion result not read
					v.addErrorAtPosition(stmt.Rhs[0].Pos(), nil)
				} else if id, ok := stmt.Lhs[1].(*ast.Ident); ok && v.blank && id.Name == "_" {
					// assertion result ignored
					v.addErrorAtPosition(id.NamePos, nil)
				}
			}
		} else {
			fmt.Println("czx@@@ assign20")
			// multiple value on rhs; in this case a call can't return
			// multiple values. Assume len(stmt.Lhs) == len(stmt.Rhs)
			for i := 0; i < len(stmt.Lhs); i++ {
				if id, ok := stmt.Lhs[i].(*ast.Ident); ok {
					if call, ok := stmt.Rhs[i].(*ast.CallExpr); ok {
						if !v.blank {
							continue
						}
						if v.ignoreCall(call) {
							continue
						}
						if id.Name == "_" && v.callReturnsError(call) {
							v.addErrorAtPosition(id.NamePos, call)
						}
					} else if assert, ok := stmt.Rhs[i].(*ast.TypeAssertExpr); ok {
						if !v.asserts {
							continue
						}
						if assert.Type == nil {
							// Shouldn't happen anyway, no multi assignment in type switches
							continue
						}
						v.addErrorAtPosition(id.NamePos, nil)
					}
				}
			}
		}
	default:
	}
	return v
}

func isErrorType(t types.Type) bool {
	return types.Implements(t, errorType)
}

func reportUncheckedErrors(e *UncheckedErrors, verbose bool) []string {
	uncheckedErrorArray := make([]string, 0)
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}
	for _, uncheckedError := range e.Errors {
		pos := uncheckedError.Pos.String()

		newPos, err := filepath.Rel(wd, pos)
		if err == nil {
			pos = newPos
		}

		if verbose && uncheckedError.FuncName != "" {
			uncheckedErrorArray = append(uncheckedErrorArray, fmt.Sprintf("%s:\t%s\t%s\n", pos, uncheckedError.FuncName, uncheckedError.Line))
		} else {
			uncheckedErrorArray = append(uncheckedErrorArray, fmt.Sprintf("%s:\t%s\n", pos, uncheckedError.Line))
		}
	}
	return uncheckedErrorArray
}
