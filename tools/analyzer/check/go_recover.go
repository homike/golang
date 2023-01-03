package check

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "goroutinewithrecover",
	Doc:      "Checks that goroutine has recover in defer function",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	// 定义过滤条件, 只分析 "go func"
	nodeFilter := []ast.Node{
		//(*ast.GoStmt)(nil),
	}

	// 检查
	inspector.Preorder(nodeFilter, func(node ast.Node) {
		pass.Reportf(123, "czx@@@ Preorder:")
		gostat, ok := node.(*ast.GoStmt)
		if ok {
			var r bool
			switch gostat.Call.Fun.(type) {
			case *ast.FuncLit: //形如 go func(){}
				funcLit := gostat.Call.Fun.(*ast.FuncLit)
				r = hasRecover(funcLit.Body)
			case *ast.Ident: //形如 go goFuncWithoutRecover()
				id := gostat.Call.Fun.(*ast.Ident)
				fd, ok := id.Obj.Decl.(*ast.FuncDecl) //fd 是 goFuncWithoutRecover 定义
				if !ok {
					return
				}
				r = hasRecover(fd.Body)
			default:

			}
			if !r {
				pass.Reportf(node.Pos(), "goroutine should have recover in defer func")
			}
		} else {
			pass.Reportf(123, "czx@@@ pass:")
			expr, ok := node.(*ast.ExprStmt)
			if ok {
				handleError(expr.X, pass)
			}
		}
	})

	return nil, nil
}

func hasRecover(bs *ast.BlockStmt) bool {
	for _, blockStmt := range bs.List {
		deferStmt, ok := blockStmt.(*ast.DeferStmt) //是否包含defer 语句
		if !ok {
			return false
		}
		switch deferStmt.Call.Fun.(type) {
		case *ast.SelectorExpr:
			//判断是否defer中包含  helper.Recover()
			selectorExpr := deferStmt.Call.Fun.(*ast.SelectorExpr)
			if "Recover" == selectorExpr.Sel.Name {
				return true
			}
		case *ast.FuncLit:
			//判断是否有 defer func(){ }()
			fl := deferStmt.Call.Fun.(*ast.FuncLit)
			for i := range fl.Body.List {

				stmt := fl.Body.List[i]
				switch stmt.(type) {
				case *ast.ExprStmt:
					exprStmt := stmt.(*ast.ExprStmt)
					if isRecoverExpr(exprStmt.X) { //recover()
						return true
					}
				case *ast.IfStmt:
					is := stmt.(*ast.IfStmt) // if r:=recover();r!=nil{}
					as, ok := is.Init.(*ast.AssignStmt)
					if !ok {
						continue
					}
					if isRecoverExpr(as.Rhs[0]) {
						return true
					}
				case *ast.AssignStmt:
					as := stmt.(*ast.AssignStmt) //r=:recover
					if isRecoverExpr(as.Rhs[0]) {
						return true
					}

				}
			}
		}
	}
	return false
}

func isRecoverExpr(expr ast.Expr) bool {
	ac, ok := expr.(*ast.CallExpr) // r:=recover()
	if !ok {
		return false
	}
	id, ok := ac.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	if "recover" == id.Name {
		return true
	}
	return false
}
