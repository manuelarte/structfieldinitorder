package internal

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// StructsHolder contains all the information of the declared structs and structs initialization.
type StructsHolder struct {
	// File
	file *ast.File
	// All the struct declarations
	structsDecl map[string]*ast.TypeSpec
	// All the struct instantiation
	structsInst map[string]*StructInit
}

func NewStructsHolder() *StructsHolder {
	structsDecl := make(map[string]*ast.TypeSpec)
	structsInst := make(map[string]*StructInit)
	return &StructsHolder{
		structsDecl: structsDecl,
		structsInst: structsInst,
	}
}

func (sh *StructsHolder) SetFile(f *ast.File) {
	sh.file = f
}

func (sh *StructsHolder) AddTypeSpec(tp *ast.TypeSpec) {
	sh.structsDecl[tp.Name.Name] = tp
}

func (sh *StructsHolder) AddCompositeLit(cl *ast.CompositeLit) {
	if si, ok := NewStructInit(cl); ok {
		sh.structsInst[si.GetIdent().Name] = si
	}
}

func (sh *StructsHolder) Analyze(pass *analysis.Pass) {
	// TODO(manuelarte): to be done
}
