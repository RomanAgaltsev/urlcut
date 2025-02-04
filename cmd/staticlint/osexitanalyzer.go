package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer - анализатор вызовов os.Exit в функции main пакета main.
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitanalyzer",
	Doc:  "check for os.Exit usage in main.main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	osExitUsed := func(x *ast.SelectorExpr) bool {
		ident, ok := x.X.(*ast.Ident)
		if !ok {
			return false
		}
		if ident.Name == "os" && x.Sel.Name == "Exit" {
			pass.Reportf(ident.NamePos, "os.Exit used in main function of main package")
			return true
		}
		return false
	}

	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}

		if strings.Contains(pass.Fset.Position(file.Pos()).Filename, "go-build") {
			continue
		}

		var isFuncMain bool

		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					isFuncMain = true
				}
			case *ast.CallExpr:
				switch y := x.Fun.(type) {
				case *ast.SelectorExpr:
					if isFuncMain && osExitUsed(y) {
						return false
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
