package internal

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// StructUniqueIdentifierKey represents how to uniquely identify a Struct.
type StructUniqueIdentifierKey struct {
	// Pkg the package the node is declared
	Pkg string
	// Name the package the node is declared
	Name string
}

func (uk StructUniqueIdentifierKey) FullImportPath(modName string) string {
	if modName == "" {
		return uk.Pkg
	}

	return fmt.Sprintf("%s/%s", modName, uk.Pkg)
}

type (
	StructSpecs struct {
		*ast.TypeSpec

		UniqueKey StructUniqueIdentifierKey
		mod       string
	}
)

func NewStructSpecs(pass *analysis.Pass, ts *ast.TypeSpec) (*StructSpecs, bool) {
	if _, ok := ts.Type.(*ast.StructType); ok {
		return &StructSpecs{
			TypeSpec: ts,
			UniqueKey: StructUniqueIdentifierKey{
				Pkg:  pass.Pkg.Path(),
				Name: ts.Name.Name,
			},
			mod: pass.Module.Path,
		}, true
	}

	return nil, false
}

// FullImportPath returns the full import path for this struct declaration.
func (ss *StructSpecs) FullImportPath() string {
	return ss.UniqueKey.FullImportPath(ss.mod)
}

func (ss *StructSpecs) GetFieldNames() []string {
	names := make([]string, 0)

	fields := ss.getStructType().Fields.List
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

func (ss *StructSpecs) getStructType() *ast.StructType {
	//nolint:errcheck // already checked
	return ss.Type.(*ast.StructType)
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
