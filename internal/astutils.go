package internal

import "go/ast"

type StructInit ast.CompositeLit

func NewStructInit(cl *ast.CompositeLit) (*StructInit, bool) {
	if _, ok := cl.Type.(*ast.Ident); ok {
		return (*StructInit)(cl), true
	}
	return nil, false
}

func (si *StructInit) GetIdent() *ast.Ident {
	return si.Type.(*ast.Ident)
}
