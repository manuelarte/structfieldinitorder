package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// StructsHolder contains all the information of the declared structs and structs initialization.
type StructsHolder struct {
	// Module name
	module string
	// File
	file *ast.File
	// imports used by the file indexed by alias.
	aliasImports map[string]string
	// All the struct declarations.
	structsDecl map[uniqueIdentifierStructKey]*ast.TypeSpec
	// All the struct instantiation.
	structsInst map[uniqueIdentifierStructKey]*StructInit
}

func NewStructsHolder(module *analysis.Module) *StructsHolder {
	aliasImports := make(map[string]string)
	structsDecl := make(map[uniqueIdentifierStructKey]*ast.TypeSpec)
	structsInst := make(map[uniqueIdentifierStructKey]*StructInit)
	return &StructsHolder{
		module:       module.Path,
		aliasImports: aliasImports,
		structsDecl:  structsDecl,
		structsInst:  structsInst,
	}
}

func (sh *StructsHolder) SetFile(f *ast.File) {
	clear(sh.aliasImports)
	sh.file = f
	sh.aliasImports[f.Name.Name] = fmt.Sprintf("%s/%s", sh.module, sh.file.Name.Name)
}

func (sh *StructsHolder) AddImportSpec(is *ast.ImportSpec) {
	// TODO(manuelarte): check for alias, and if not, get the latest name after /
	if is.Path != nil && is.Path.Kind == token.STRING {
		// remove ""
		alias := is.Path.Value[1 : len(is.Path.Value)-1]
		indexDel := strings.LastIndex(alias, "/")
		if indexDel != -1 {
			alias = alias[indexDel+1:]
		}
		if is.Name != nil {
			alias = is.Name.Name
		}
		sh.aliasImports[alias] = is.Path.Value
	}
}

func (sh *StructsHolder) AddTypeSpec(tp *ast.TypeSpec) {
	if _, ok := tp.Type.(*ast.StructType); ok {
		importPath := fmt.Sprintf("%s/%s", sh.module, sh.file.Name)
		sh.structsDecl[uniqueIdentifierStructKey{importPath: importPath, structName: tp.Name.Name}] = tp
	}
}

func (sh *StructsHolder) AddCompositeLit(cl *ast.CompositeLit) {
	if si, ok := NewStructInit(sh.file.Name.Name, cl); ok {
		// TODO(manuelarte): substitute the pkg alias with the unique identifier
		importPath := sh.aliasImports[si.key.packageAlias]
		sh.structsInst[uniqueIdentifierStructKey{
			importPath: importPath,
			structName: si.key.name,
		}] = si
	}
}

func (sh *StructsHolder) Analyze(pass *analysis.Pass) {
	// TODO(manuelarte): to be done
}

type uniqueIdentifierStructKey struct {
	importPath string
	structName string
}
