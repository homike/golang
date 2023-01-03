package check

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// error必须被处理,
// 1. return
// 2. logger打印

// _, err := func()
func handleError(expr ast.Expr, pass *analysis.Pass) bool {
	fmt.Println("handleError")
	isErrorExpr(expr, pass)
	return true
}

/*
	for _, blockStmt := range bs.List {
		deferStmt, ok := blockStmt.(*ast.Expr) //是否包含defer 语句
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
*/

func isErrorExpr(expr ast.Expr, pass *analysis.Pass) bool {
	ac, ok := expr.(*ast.BinaryExpr) // r:=recover()
	if !ok {
		pass.Reportf(0, "11name")
		return false
	}
	id, ok := ac.X.(*ast.Ident)
	if !ok {
		pass.Reportf(0, "22name:")
		return false
	}

	pass.Reportf(0, "name:")
	if "recover" == id.Name {
		return true
	}
	return false
}

/*
func isError(v ast.Expr, info *types.Info) bool {
	if intf, ok := info.TypeOf(v).Underlying().(*types.Interface); ok {
		return intf.NumMethods() == 1 && intf.Method(0).FullName() == "(error).Error"
	}
	return false
}
*/
