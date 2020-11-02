package golinters

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewCheckTimeNow() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"bannedfunc",
		"Checks that cannot use func",
		[]*analysis.Analyzer{Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

var Analyzer = &analysis.Analyzer{
	Name:     "time",
	Doc:      "检查是否使用 time.Now()",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	astf := func(node ast.Node) bool {
		selector, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := selector.X.(*ast.Ident)
		if !ok {
			return true
		}
		if ident.Name != "time" {
			return true
		}
		sel := selector.Sel
		if sel.Name != "Now" {
			return true
		}
		pass.ReportRangef(node, "不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()")
		return true
	}
	for _, f := range pass.Files {
		ast.Inspect(f, astf)
	}
	return nil, nil
}
