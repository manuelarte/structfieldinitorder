package analyzer

import (
	"go/ast"
	"go/types"
	"maps"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/manuelarte/structfieldinitorder/internal"
)

type structFieldInitOrder struct {
	mu sync.RWMutex
	// all the struct specs of the codebase.
	structsSpec map[internal.UniqueIdentifierStructKey]*internal.StructSpec
	// analysis.Pass and struct instantiated indexed by package.
	stateIndexByPkg map[*types.Package]state
}

func NewAnalyzer() *analysis.Analyzer {
	structsSpec := make(map[internal.UniqueIdentifierStructKey]*internal.StructSpec)
	st := make(map[*types.Package]state)
	s := structFieldInitOrder{
		structsSpec:     structsSpec,
		stateIndexByPkg: st,
	}

	return &analysis.Analyzer{
		Name:     "structfieldinitorder",
		Doc:      "This linter checks whether when a struct is instantiated, the fields order follows the same order as in the struct declaration.", //nolint:lll // url
		URL:      "https://github.com/manuelarte/structfieldinitorder",
		Run:      s.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func (s *structFieldInitOrder) run(pass *analysis.Pass) (any, error) {
	insp, found := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !found {
		//nolint:nilnil // impossible case.
		return nil, nil
	}

	fh := internal.NewStructsHolder(pass.Pkg.Path())

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.ImportSpec)(nil),
		(*ast.TypeSpec)(nil),
		(*ast.CompositeLit)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.File:
			fh.SetFile(node)
		case *ast.ImportSpec:
			fh.AddImportSpec(node)
		case *ast.TypeSpec:
			fh.AddTypeSpec(node)
		case *ast.CompositeLit:
			fh.AddCompositeLit(node)
		}
	})

	s.mu.Lock()
	defer s.mu.Unlock()
	maps.Copy(s.structsSpec, fh.StructsSpec())
	if len(fh.StructInst()) > 0 {
		if _, ok := s.stateIndexByPkg[pass.Pkg]; !ok {
			s.stateIndexByPkg[pass.Pkg] = state{
				pass:        pass,
				structInsts: fh.StructInst(),
			}
		} else {
			st := s.stateIndexByPkg[pass.Pkg]
			st.copy(fh.StructInst())
		}
	}

	fh.CheckDotImports()

	s.analyze()

	//nolint:nilnil //any, error
	return nil, nil
}

func (s *structFieldInitOrder) analyze() {
	for _, st := range s.stateIndexByPkg {
		for key, structInsts := range st.structInsts {
			if structSpec, ok := s.structsSpec[key]; ok {
				for _, structInst := range structInsts {
					internal.ReportIfStructFieldsNotInOrder(st.pass, structSpec, structInst)
				}
				delete(st.structInsts, key)
			}
		}
	}
}

type state struct {
	pass        *analysis.Pass
	structInsts map[internal.UniqueIdentifierStructKey][]*internal.StructInst
}

func (s *state) copy(c map[internal.UniqueIdentifierStructKey][]*internal.StructInst) {
	maps.Copy(s.structInsts, c)
}
