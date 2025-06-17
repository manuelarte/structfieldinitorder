package internal

import (
	"go/ast"
	"go/token"
	"strings"
)

type IStructInst interface {
	ast.Node

	GetFieldNames() []string
	x()
}

func NewIStructInst(pkgPath string, importsSpec []*ast.ImportSpec, cl *ast.CompositeLit) (IStructInst, bool) {
	if len(cl.Elts) == 0 {
		return nil, false
	}
	for _, elt := range cl.Elts {
		if _, ok := elt.(*ast.KeyValueExpr); !ok {
			return nil, false
		}
	}
	switch cl.Type.(type) {
	case *ast.Ident:
		// TODO(manuelarte): case in which either the struct is declared in the same package or it's using a dot import.
		return newStructInstInSamePkgStructDecl(pkgPath, cl)

	case *ast.SelectorExpr:
		// this means the struct is declared in a different package
		return newStructInstWithAlias(importsSpec, cl)
	default:
		return nil, false
	}
}

type StructInstWithAlias struct {
	*ast.CompositeLit

	// Name of the struct
	Name string
	// ImportSpec that matches the pkg alias.
	ImportSpec *ast.ImportSpec
}

func newStructInstWithAlias(
	importsSpec []*ast.ImportSpec,
	cl *ast.CompositeLit,
) (*StructInstWithAlias, bool) {
	//nolint: nestif // I may remove the defensive programming type check of ast.SelectorExpr
	if selectorExpr, ok := cl.Type.(*ast.SelectorExpr); ok {
		if xIdent, isXIdent := selectorExpr.X.(*ast.Ident); isXIdent {
			pkgAlias := xIdent.Name
			for _, is := range importsSpec {
				if importAlias, found := getPkgImportAlias(is); found {
					if importAlias == pkgAlias {
						return &StructInstWithAlias{
							CompositeLit: cl,
							Name:         selectorExpr.Sel.Name,
							ImportSpec:   is,
						}, true
					}
				}
			}
		}
	}
	return nil, false
}

// GetStructUniqueIdentifierKey returns the unique identifier where the struct was declared.
func (si *StructInstWithAlias) GetStructUniqueIdentifierKey() StructUniqueIdentifierKey {
	importPath, _ := getPkgImportPath(si.ImportSpec)
	return StructUniqueIdentifierKey{
		Pkg:  importPath,
		Name: si.Name,
	}
}

func (si *StructInstWithAlias) GetFieldNames() []string {
	names := make([]string, 0)
	kvs := si.getKeyValueExpr()
	for _, kv := range kvs {
		if ident, ok := kv.Key.(*ast.Ident); ok {
			names = append(names, ident.Name)
		}
	}
	return names
}

func (si *StructInstWithAlias) getKeyValueExpr() []*ast.KeyValueExpr {
	kv := make([]*ast.KeyValueExpr, len(si.Elts))
	for i, elt := range si.Elts {
		//nolint:errcheck // already checked
		kv[i] = elt.(*ast.KeyValueExpr)
	}
	return kv
}

func (*StructInstWithAlias) x() {}

type StructInstInSamePkgStructDecl struct {
	*ast.CompositeLit

	// Name of the struct
	Name string
	// ImportPath of the pkg
	ImportPath string
}

func newStructInstInSamePkgStructDecl(
	pkgPath string,
	cl *ast.CompositeLit,
) (*StructInstInSamePkgStructDecl, bool) {
	if ident, ok := cl.Type.(*ast.Ident); ok {
		return &StructInstInSamePkgStructDecl{
			CompositeLit: cl,
			Name:         ident.Name,
			ImportPath:   pkgPath,
		}, true
	}
	return nil, false
}

// GetStructUniqueIdentifierKey returns the unique identifier where the struct was declared.
func (si *StructInstInSamePkgStructDecl) GetStructUniqueIdentifierKey() StructUniqueIdentifierKey {
	return StructUniqueIdentifierKey{
		Pkg:  si.ImportPath,
		Name: si.Name,
	}
}

func (si *StructInstInSamePkgStructDecl) GetFieldNames() []string {
	names := make([]string, 0)
	kvs := si.getKeyValueExpr()
	for _, kv := range kvs {
		if ident, ok := kv.Key.(*ast.Ident); ok {
			names = append(names, ident.Name)
		}
	}
	return names
}

func (si *StructInstInSamePkgStructDecl) getKeyValueExpr() []*ast.KeyValueExpr {
	kv := make([]*ast.KeyValueExpr, len(si.Elts))
	for i, elt := range si.Elts {
		//nolint:errcheck // already checked
		kv[i] = elt.(*ast.KeyValueExpr)
	}
	return kv
}

func (*StructInstInSamePkgStructDecl) x() {}

func getPkgImportAlias(is *ast.ImportSpec) (string, bool) {
	if is.Path != nil && is.Path.Kind == token.STRING {
		// remove ""
		alias := is.Path.Value[1 : len(is.Path.Value)-1]
		indexDel := strings.LastIndex(alias, "/")
		if indexDel != -1 {
			alias = alias[indexDel+1:]
		}
		if is.Name != nil && is.Name.Name != "." {
			alias = is.Name.Name
		}
		return alias, true
	}
	return "", false
}

func getPkgImportPath(is *ast.ImportSpec) (string, bool) {
	if is.Path != nil && is.Path.Kind == token.STRING {
		// remove ""
		return is.Path.Value[1 : len(is.Path.Value)-1], true
	}
	return "", false
}
