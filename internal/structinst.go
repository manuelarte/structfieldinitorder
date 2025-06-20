package internal

import (
	"go/ast"
	"go/token"
	"strings"
)

type StructInst interface {
	GetCompositeLit() *ast.CompositeLit
	GetKeyValueExpr() []*ast.KeyValueExpr
	x()
}

func NewStructInst(pkgPath string, importsSpec []*ast.ImportSpec, cl *ast.CompositeLit) (StructInst, bool) {
	bsi, ok := newBaseStructInst(cl)
	if !ok {
		return nil, false
	}
	switch cl.Type.(type) {
	case *ast.Ident:
		// this means the struct is either declared in the same package or there are dot imports.
		dotImports := make([]*ast.ImportSpec, 0, len(importsSpec))
		for _, spec := range importsSpec {
			if spec.Name != nil && spec.Name.Name == "." {
				dotImports = append(dotImports, spec)
			}
		}
		if len(dotImports) == 0 {
			return newStructInstInSamePkgStructDecl(pkgPath, bsi)
		}
		return newStructInstWithDotImports(pkgPath, dotImports, bsi)
	case *ast.SelectorExpr:
		// this means the struct is declared in a different package
		return newStructInstWithAlias(importsSpec, bsi)
	default:
		return nil, false
	}
}

type baseStructInst struct {
	*ast.CompositeLit
}

func newBaseStructInst(cl *ast.CompositeLit) (*baseStructInst, bool) {
	if len(cl.Elts) == 0 {
		return nil, false
	}
	for _, elt := range cl.Elts {
		if _, ok := elt.(*ast.KeyValueExpr); !ok {
			return nil, false
		}
	}

	return &baseStructInst{
		cl,
	}, true
}

func (si *baseStructInst) GetCompositeLit() *ast.CompositeLit {
	return si.CompositeLit
}

func (si *baseStructInst) GetKeyValueExpr() []*ast.KeyValueExpr {
	kvs, _ := transform(si.Elts, func(i ast.Expr) (*ast.KeyValueExpr, error) {
		//nolint:errcheck // checked before
		return i.(*ast.KeyValueExpr), nil
	})
	return kvs
}

func (*baseStructInst) x() {}

type (
	StructInstWithAlias struct {
		*baseStructInst

		// Name of the struct.
		Name string
		// ImportSpec that matches the pkg alias.
		ImportSpec *ast.ImportSpec
	}

	StructInstInSamePkgStructDecl struct {
		*baseStructInst

		// Name of the struct.
		Name string
		// ImportPath of the pkg.
		ImportPath string
	}

	StructInstWithDotImports struct {
		*baseStructInst

		// Name of the struct.
		Name string
		// ImportPath of the pkg.
		ImportPath string
		// Dot Import Specs, contains all the dot imports of the file where the struct was instantiated.
		DotImports []*ast.ImportSpec
	}
)

func newStructInstWithAlias(
	importsSpec []*ast.ImportSpec,
	si *baseStructInst,
) (*StructInstWithAlias, bool) {
	//nolint: nestif // I may remove the defensive programming type check of ast.SelectorExpr
	if selectorExpr, ok := si.Type.(*ast.SelectorExpr); ok {
		if xIdent, isXIdent := selectorExpr.X.(*ast.Ident); isXIdent {
			pkgAlias := xIdent.Name
			for _, is := range importsSpec {
				if importAlias, found := getPkgImportAlias(is); found {
					if importAlias == pkgAlias {
						return &StructInstWithAlias{
							baseStructInst: si,
							Name:           selectorExpr.Sel.Name,
							ImportSpec:     is,
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

func newStructInstInSamePkgStructDecl(
	pkgPath string,
	si *baseStructInst,
) (*StructInstInSamePkgStructDecl, bool) {
	if ident, ok := si.Type.(*ast.Ident); ok {
		return &StructInstInSamePkgStructDecl{
			baseStructInst: si,
			Name:           ident.Name,
			ImportPath:     pkgPath,
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

func newStructInstWithDotImports(
	pkgPath string,
	dotImports []*ast.ImportSpec,
	si *baseStructInst,
) (*StructInstWithDotImports, bool) {
	if ident, ok := si.Type.(*ast.Ident); ok {
		return &StructInstWithDotImports{
			baseStructInst: si,
			Name:           ident.Name,
			ImportPath:     pkgPath,
			DotImports:     dotImports,
		}, true
	}
	return nil, false
}

func (si *StructInstWithDotImports) GetMatchingStructSpecs(
	structSpecsIndexed map[StructUniqueIdentifierKey]*StructSpecs,
) (*StructSpecs, bool) {
	// check pkgPath first
	if ss, ok := structSpecsIndexed[StructUniqueIdentifierKey{
		Pkg:  si.ImportPath,
		Name: si.Name,
	}]; ok {
		return ss, true
	}
	for _, importSpec := range si.DotImports {
		importPath, foundImportPath := getPkgImportPath(importSpec)
		if !foundImportPath {
			continue
		}
		key := StructUniqueIdentifierKey{Pkg: importPath, Name: si.Name}
		if ss, ok := structSpecsIndexed[key]; ok {
			return ss, true
		}
	}
	return &StructSpecs{}, false
}

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
