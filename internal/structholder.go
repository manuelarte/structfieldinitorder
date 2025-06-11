package internal

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// StructsHolder contains all the information of the declared structs and structs initialization.
type StructsHolder struct {
	// All the struct declarations
	structs map[string]*ast.TypeSpec
}

func NewStructsHolder() *StructsHolder {
	structs := make(map[string]*ast.TypeSpec)
	return &StructsHolder{
		structs: structs,
	}
}

func (sh *StructsHolder) AddTypeSpec(tp *ast.TypeSpec) {
	sh.structs[tp.Name.Name] = tp
}

func (sh *StructsHolder) Analyze(pass *analysis.Pass) {
	// TODO(manuelarte): to be done
}
