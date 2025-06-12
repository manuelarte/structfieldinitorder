package internal

import "go/ast"

type StructInit struct {
	key aliasIdentifierStructKey
	*ast.CompositeLit
}

func NewStructInit(pkgName string, cl *ast.CompositeLit) (*StructInit, bool) {
	if key, ok := getKey(pkgName, cl); ok {
		return &StructInit{
			key:          key,
			CompositeLit: cl,
		}, true
	}
	return nil, false
}

func getKey(pkgName string, cl *ast.CompositeLit) (aliasIdentifierStructKey, bool) {
	switch node := cl.Type.(type) {
	case *ast.Ident:
		// this means the struct is declared in the same package
		return aliasIdentifierStructKey{
			packageAlias: pkgName,
			name:         node.Name,
		}, true
	case *ast.SelectorExpr:
		if xIdent, isXIdent := node.X.(*ast.Ident); isXIdent {
			return aliasIdentifierStructKey{
				packageAlias: xIdent.Name,
				name:         node.Sel.Name,
			}, true
		}
		return aliasIdentifierStructKey{}, false
	default:
		return aliasIdentifierStructKey{}, false
	}
}

// key to identify a struct, e.g.
// e.g.
// packageAlias: internal
// name: structIdentifierKey
type aliasIdentifierStructKey struct {
	packageAlias string
	name         string
}
