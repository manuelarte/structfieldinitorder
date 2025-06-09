package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/manuelarte/structfieldinitorder/internal"
)

func NewAnalyzer() *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name:     "structfieldinitorder",
		Doc:      "this linter checks whether when a struct is instantiated, the fields order follows the same order as in the struct declaration.",
		URL:      "https://github.com/manuelarte/structfieldinitorder",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	return a
}

func run(pass *analysis.Pass) (any, error) {
	insp, found := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !found {
		//nolint:nilnil // impossible case.
		return nil, nil
	}

	sh := internal.NewStructHolder()

	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
		//TODO(manuelarte): missing struct initialization
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.TypeSpec:
			sh.AddTypeSpec(node)
		}
	})

	sh.Analyze(pass)

	//nolint:nilnil //any, error
	return nil, nil
}
