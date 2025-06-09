package internal

import (
	"cmp"
	"go/ast"
	"slices"

	"golang.org/x/tools/go/analysis"
)

// StructsHolder contains all the information of the declared structs and structs initialization.
type StructsHolder struct {
	// The struct declaration
	Struct *ast.TypeSpec
}

func NewStructsHolder() *StructsHolder {
	return &StructsHolder{}
}

func (sh *StructsHolder) AddTypeSpec(tp *ast.TypeSpec) {
	// TODO(manuelarte): to be done
}

func (sh *StructsHolder) Analyze(pass *analysis.Pass) {
	// TODO(manuelarte): to be done
}

