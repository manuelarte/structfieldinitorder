package internal

import (
	"go/ast"
	"go/token"
	"maps"
	"strings"
)

// StructsHolder contains all the information of the declared structs and structs initialization.
type StructsHolder struct {
	pkgPath string
	file    *ast.File
	// imports used by the file, indexed by alias.
	aliasImports map[string]string
	// all the struct specs of the file.
	structsSpec map[UniqueIdentifierStructKey]*StructSpec
	// all the struct instantiation of the file.
	structsInst map[UniqueIdentifierStructKey][]*StructInst
}

func NewStructsHolder(pkgPath string) *StructsHolder {
	aliasImports := make(map[string]string)
	structsSpec := make(map[UniqueIdentifierStructKey]*StructSpec)
	structsInst := make(map[UniqueIdentifierStructKey][]*StructInst)
	return &StructsHolder{
		pkgPath:      pkgPath,
		aliasImports: aliasImports,
		structsSpec:  structsSpec,
		structsInst:  structsInst,
	}
}

func (sh *StructsHolder) StructsSpec() map[UniqueIdentifierStructKey]*StructSpec {
	return maps.Clone(sh.structsSpec)
}

func (sh *StructsHolder) StructInst() map[UniqueIdentifierStructKey][]*StructInst {
	return maps.Clone(sh.structsInst)
}

func (sh *StructsHolder) SetFile(f *ast.File) {
	clear(sh.aliasImports)
	sh.file = f
	sh.aliasImports[f.Name.Name] = sh.pkgPath
}

func (sh *StructsHolder) AddImportSpec(is *ast.ImportSpec) {
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
		sh.aliasImports[alias] = is.Path.Value[1 : len(is.Path.Value)-1]
	}
}

func (sh *StructsHolder) AddTypeSpec(tp *ast.TypeSpec) {
	if ss, ok := NewStructSpec(sh.file.Name.Name, tp); ok {
		importPath := sh.pkgPath
		sh.structsSpec[UniqueIdentifierStructKey{ImportPath: importPath, StructName: ss.key.name}] = ss
	}
}

func (sh *StructsHolder) AddCompositeLit(cl *ast.CompositeLit) {
	if si, ok := NewStructInit(sh.file.Name.Name, cl); ok {
		if importPath, exists := sh.aliasImports[si.key.pkgAlias]; exists {
			key := UniqueIdentifierStructKey{
				ImportPath: importPath,
				StructName: si.key.name,
			}
			sh.structsInst[key] = append(sh.structsInst[key], si)
		}
	}
}

type UniqueIdentifierStructKey struct {
	ImportPath string
	StructName string
}
