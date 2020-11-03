package golinters

import (
	"go/ast"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Configuration represents go-header linter setup parameters
type Configuration struct {
	// Values is map of values. Supports two types 'const` and `regexp`. Values can be used recursively.
	Values map[string]map[string]string `yaml:"values"`
}

var Analyzer = &analysis.Analyzer{
	Name:     "time",
	Doc:      "检查配置里列出的函数调用",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func checkFirstLetter(s string) (string, bool) {
	if len(s) < 2 || s[0] != '~' || s[1] < 97 || s[1] > 122 {
		return "", false
	}
	return string(s[1]-32) + s[2:], true
}
func NewCheckTimeNow() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"bannedfunc",
		"Checks that cannot use func",
		[]*analysis.Analyzer{Analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		cfg := lintCtx.Cfg.LintersSettings.BannedFunc
		for k, v := range cfg.Values {
			for itemK, item := range v {
				if s, ok := checkFirstLetter(itemK); ok {
					v[s] = item
					delete(v, itemK)
				}
			}
			if s, ok := checkFirstLetter(k); ok {
				delete(cfg.Values, k)
				k = s
			}
			cfg.Values[k] = v
		}
		c := &Configuration{
			Values: cfg.Values,
		}

		Analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			astf := func(node ast.Node) bool {
				selector, ok := node.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				ident, ok := selector.X.(*ast.Ident)
				if !ok {
					return true
				}

				m, ok := c.Values[ident.Name]
				if !ok {
					return true
				}

				sel := selector.Sel
				value, ok := m[sel.Name]
				if !ok {
					return true
				}
				pass.Reportf(node.Pos(), value)
				return true
			}
			for _, f := range pass.Files {
				ast.Inspect(f, astf)
			}
			return nil, nil
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
