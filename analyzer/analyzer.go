package analyzer

import (
	"context"
	"fmt"
	"go/ast"

	"github.com/ldez/grignotin/gomod"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/manuelarte/structfieldinitorder/internal"
)

func NewAnalyzer() *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name:     "structfieldinitorder",
		Doc:      "this linter checks whether when a struct is instantiated, the fields order follows the same order as in the struct declaration.", //nolint:lll // url
		URL:      "https://github.com/manuelarte/structfieldinitorder",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	return a
}

func run(pass *analysis.Pass) (any, error) {
	info, err := gomod.GetModuleInfo(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(info)

	insp, found := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !found {
		//nolint:nilnil // impossible case.
		return nil, nil
	}

	sh := internal.NewStructsHolder()

	// TODO(manuelarte): I think this does not work in this linter because I need the package.
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.TypeSpec)(nil),
		(*ast.CompositeLit)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.File:
			sh.SetFile(node)
		case *ast.TypeSpec:
			sh.AddTypeSpec(node)
		case *ast.CompositeLit:
			sh.AddCompositeLit(node)
		}
	})

	sh.Analyze(pass)

	//nolint:nilnil //any, error
	return nil, nil
}
