// Package noexit содержит линтер для проверки использование прямого вызова os.Exit
// в функции main пакета main.
package noexit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noexit",
	Doc:  "проверка использования прямого вызова os.Exit в функции main пакета main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		// Обрабатываем только go файлы, также исключаем из проверки тесты.
		if !strings.HasSuffix(filename, ".go") || strings.HasSuffix(filename, "_test.go") {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				return isMainFunc(x)
			case *ast.ExprStmt:
				validateExprStmt(x, pass)
			case *ast.GoStmt:
				validateGoStmt(x, pass)
			case *ast.DeferStmt:
				validateDeferStmt(x, pass)
			case *ast.FuncLit:
				return false
			}
			return true
		})
	}

	return nil, nil //nolint: nilnil
}

// isMainFunc возвращает true если нода является функцией main.
func isMainFunc(node *ast.FuncDecl) bool {
	if node.Name == nil || node.Name.Obj == nil {
		return false
	}

	fn := node.Name.Obj

	return fn.Kind == ast.Fun && fn.Name == "main"
}

func validateCallExpr(node ast.Node, pass *analysis.Pass) {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return
	}

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	pkg, ok := selector.X.(*ast.Ident)
	if !ok {
		return
	}

	if pkg.Name == "os" && selector.Sel.Name == "Exit" {
		pass.Reportf(node.Pos(), "запрещено использовать прямой вызов os.Exit в функции main пакета main")
		return
	}
}

func validateExprStmt(node *ast.ExprStmt, pass *analysis.Pass) {
	validateCallExpr(node.X, pass)
}

func validateGoStmt(node *ast.GoStmt, pass *analysis.Pass) {
	if node.Call != nil {
		validateCallExpr(node.Call, pass)
	}
}

func validateDeferStmt(node *ast.DeferStmt, pass *analysis.Pass) {
	if node.Call != nil {
		validateCallExpr(node.Call, pass)
	}
}
