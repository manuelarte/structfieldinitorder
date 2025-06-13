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
	// module name.
	module string
	file   *ast.File
	// imports used by the file, indexed by alias.
	aliasImports map[string]string
	// all the struct specs that will be checked.
	structsSpec map[uniqueIdentifierStructKey]*StructSpec
	// all the struct instantiation that will be checked.
	structsInst map[uniqueIdentifierStructKey][]*StructInst
}

func NewStructsHolder(module *analysis.Module) *StructsHolder {
	aliasImports := make(map[string]string)
	structsSpec := make(map[uniqueIdentifierStructKey]*StructSpec)
	structsInst := make(map[uniqueIdentifierStructKey][]*StructInst)
	return &StructsHolder{
		module:       module.Path,
		aliasImports: aliasImports,
		structsSpec:  structsSpec,
		structsInst:  structsInst,
	}
}

func (sh *StructsHolder) SetFile(f *ast.File) {
	clear(sh.aliasImports)
	sh.file = f
	sh.aliasImports[f.Name.Name] = fmt.Sprintf("%s/%s", sh.module, sh.file.Name.Name)
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
		sh.aliasImports[alias] = is.Path.Value
	}
}

func (sh *StructsHolder) AddTypeSpec(tp *ast.TypeSpec) {
	if ss, ok := NewStructSpec(sh.file.Name.Name, tp); ok {
		importPath := fmt.Sprintf("%s/%s", sh.module, sh.file.Name)
		sh.structsSpec[uniqueIdentifierStructKey{importPath: importPath, structName: ss.key.name}] = ss
	}
}

func (sh *StructsHolder) AddCompositeLit(cl *ast.CompositeLit) {
	if si, ok := NewStructInit(sh.file.Name.Name, cl); ok {
		if importPath, exists := sh.aliasImports[si.key.pkgAlias]; exists {
			key := uniqueIdentifierStructKey{
				importPath: importPath,
				structName: si.key.name,
			}
			sh.structsInst[key] = append(sh.structsInst[key], si)
		}
	}
}

func (sh *StructsHolder) Analyze(pass *analysis.Pass) {
	for key, structInsts := range sh.structsInst {
		if structSpec, ok := sh.structsSpec[key]; ok {
			for _, structInst := range structInsts {
				reportIfStructFieldsNotInOrder(pass, structSpec, structInst)
			}
		}
	}
}

type uniqueIdentifierStructKey struct {
	importPath string
	structName string
}
