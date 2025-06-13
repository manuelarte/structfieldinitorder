package internal

import (
	"go/ast"
)

// StructInst holds data of a struct initialization.
type StructInst struct {
	*ast.CompositeLit

	key aliasIdentifierStructKey
}

func NewStructInit(pkgName string, cl *ast.CompositeLit) (*StructInst, bool) {
	if len(cl.Elts) == 0 {
		return nil, false
	}
	for _, elt := range cl.Elts {
		if _, ok := elt.(*ast.KeyValueExpr); !ok {
			return nil, false
		}
	}
	if key, ok := getKey(pkgName, cl); ok {
		return &StructInst{
			CompositeLit: cl,
			key:          key,
		}, true
	}
	return nil, false
}

func (si *StructInst) GetKeyValueExpr() []*ast.KeyValueExpr {
	kv := make([]*ast.KeyValueExpr, len(si.Elts))
	for i, elt := range si.Elts {
		//nolint:errcheck // already checked
		kv[i] = elt.(*ast.KeyValueExpr)
	}
	return kv
}

func (si *StructInst) GetFieldNames() []string {
	names := make([]string, 0)
	kvs := si.GetKeyValueExpr()
	for _, kv := range kvs {
		if ident, ok := kv.Key.(*ast.Ident); ok {
			names = append(names, ident.Name)
		}
	}
	return names
}

func getKey(pkgName string, cl *ast.CompositeLit) (aliasIdentifierStructKey, bool) {
	switch node := cl.Type.(type) {
	case *ast.Ident:
		// this means the struct is declared in the same package
		return aliasIdentifierStructKey{
			pkgAlias: pkgName,
			name:     node.Name,
		}, true
	case *ast.SelectorExpr:
		if xIdent, isXIdent := node.X.(*ast.Ident); isXIdent {
			return aliasIdentifierStructKey{
				pkgAlias: xIdent.Name,
				name:     node.Sel.Name,
			}, true
		}
		return aliasIdentifierStructKey{}, false
	default:
		return aliasIdentifierStructKey{}, false
	}
}

type StructSpec struct {
	*ast.TypeSpec

	key aliasIdentifierStructKey
}

func NewStructSpec(pkgName string, ts *ast.TypeSpec) (*StructSpec, bool) {
	if _, ok := ts.Type.(*ast.StructType); ok {
		return &StructSpec{
			TypeSpec: ts,
			key: aliasIdentifierStructKey{
				pkgAlias: pkgName,
				name:     ts.Name.Name,
			},
		}, true
	}
	return nil, false
}

func (si *StructSpec) GetStructType() *ast.StructType {
	//nolint:errcheck // already checked
	return si.Type.(*ast.StructType)
}

func (si *StructSpec) GetFieldNames() []string {
	names := make([]string, 0)
	fields := si.GetStructType().Fields.List
	for _, field := range fields {
		if field.Names != nil {
			for _, name := range field.Names {
				names = append(names, name.Name)
			}
		} else {
			// embedded type
			if name, ok := getFieldNameFromType(field.Type); ok {
				names = append(names, name)
			}
		}
	}

	return names
}

func getFieldNameFromType(expr ast.Expr) (string, bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, true
	case *ast.StarExpr:
		return getFieldNameFromType(t.X)
	case *ast.SelectorExpr:
		return getFieldNameFromType(t.Sel)
	default:
		return "", false
	}
}

// key to identify a struct, e.g.
// e.g.
// pkgAlias: internal
// name: structIdentifierKey
type aliasIdentifierStructKey struct {
	pkgAlias string
	name     string
}
