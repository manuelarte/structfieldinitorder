package analyzer

import (
	"go/ast"
	"go/types"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/manuelarte/structfieldinitorder/internal"
)

type structFieldInitOrder struct {
	mu sync.RWMutex

	structSpecsIndexedByKey map[internal.StructUniqueIdentifierKey]*internal.StructSpecs
	structInstIndexedByPkg  map[*types.Package]*pkgStructInst
}

type pkgStructInst struct {
	pass    *analysis.Pass
	structs []internal.StructInst
}

func (s *pkgStructInst) append(si internal.StructInst) {
	s.structs = append(s.structs, si)
}

func NewAnalyzer() *analysis.Analyzer {
	s := structFieldInitOrder{
		structSpecsIndexedByKey: make(map[internal.StructUniqueIdentifierKey]*internal.StructSpecs),
		structInstIndexedByPkg:  make(map[*types.Package]*pkgStructInst),
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
	s.mu.Lock()
	defer s.mu.Unlock()

	insp, found := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !found {
		//nolint:nilnil // impossible case.
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.ImportSpec)(nil),
		(*ast.TypeSpec)(nil),
		(*ast.CompositeLit)(nil),
	}

	importsSpec := make([]*ast.ImportSpec, 0)
	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.File:
			importsSpec = nil
		case *ast.ImportSpec:
			importsSpec = append(importsSpec, node)
		case *ast.TypeSpec:
			{
				if ss, ok := internal.NewStructSpecs(pass, node); ok {
					s.structSpecsIndexedByKey[ss.UniqueKey] = ss
				}
			}

		case *ast.CompositeLit:
			if si, ok := internal.NewStructInst(pass.Pkg.Path(), importsSpec, node); ok {
				{
					// add si
					if pkgKey, pkgFound := s.structInstIndexedByPkg[pass.Pkg]; pkgFound {
						pkgKey.append(si)
					} else {
						s.structInstIndexedByPkg[pass.Pkg] = &pkgStructInst{
							pass:    pass,
							structs: []internal.StructInst{si},
						}
					}
				}
			}
		}
	})

	importsSpec = nil

	s.analyze()

	//nolint:nilnil //any, error
	return nil, nil
}

func (s *structFieldInitOrder) analyze() {
	for _, st := range s.structInstIndexedByPkg {
		notProcessedStructs := make([]internal.StructInst, 0)
		for _, structInst := range st.structs {
			var structUniqueIdentifierKey internal.StructUniqueIdentifierKey
			switch si := structInst.(type) {
			case *internal.StructInstWithAlias:
				structUniqueIdentifierKey = si.GetStructUniqueIdentifierKey()
			case *internal.StructInstInSamePkgStructDecl:
				structUniqueIdentifierKey = si.GetStructUniqueIdentifierKey()
			case *internal.StructInstWithDotImports:
				structSpecs, _ := si.GetMatchingStructSpecs(s.structSpecsIndexedByKey)
				structUniqueIdentifierKey = structSpecs.UniqueKey
			default:
				continue
			}

			if structSpec, ok := s.structSpecsIndexedByKey[structUniqueIdentifierKey]; ok {
				internal.ReportIfStructFieldsNotInOrder(st.pass, structSpec, structInst)
			} else {
				notProcessedStructs = append(notProcessedStructs, structInst)
			}
		}

		st.structs = notProcessedStructs
	}
}
