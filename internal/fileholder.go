package internal

import (
	"go/ast"
	"go/token"
	"maps"
	"strings"
)

// FileHolder contains all the information of the declared structs and structs initialization.
type FileHolder struct {
	pkgPath string
	file    *ast.File
	// imports used by the file, indexed by alias.
	aliasImports map[string]string
	// all the struct specs of the file.
	structsSpec map[UniqueIdentifierStructKey]*StructSpec
	// all the struct instantiation of the file.
	structsInst map[UniqueIdentifierStructKey][]*StructInst
}

func NewStructsHolder(pkgPath string) *FileHolder {
	aliasImports := make(map[string]string)
	structsSpec := make(map[UniqueIdentifierStructKey]*StructSpec)
	structsInst := make(map[UniqueIdentifierStructKey][]*StructInst)
	return &FileHolder{
		pkgPath:      pkgPath,
		aliasImports: aliasImports,
		structsSpec:  structsSpec,
		structsInst:  structsInst,
	}
}

func (fh *FileHolder) StructsSpec() map[UniqueIdentifierStructKey]*StructSpec {
	return maps.Clone(fh.structsSpec)
}

func (fh *FileHolder) StructInst() map[UniqueIdentifierStructKey][]*StructInst {
	return maps.Clone(fh.structsInst)
}

func (fh *FileHolder) SetFile(f *ast.File) {
	clear(fh.aliasImports)
	fh.file = f
	if f != nil {
		fh.aliasImports[f.Name.Name] = fh.pkgPath
	}
}

func (fh *FileHolder) AddImportSpec(is *ast.ImportSpec) {
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
		fh.aliasImports[alias] = is.Path.Value[1 : len(is.Path.Value)-1]
	}
}

func (fh *FileHolder) AddTypeSpec(tp *ast.TypeSpec) {
	if ss, ok := NewStructSpec(fh.file.Name.Name, tp); ok {
		importPath := fh.pkgPath
		fh.structsSpec[UniqueIdentifierStructKey{ImportPath: importPath, StructName: ss.key.name}] = ss
	}
}

func (fh *FileHolder) AddCompositeLit(cl *ast.CompositeLit) {
	if si, ok := NewStructInit(fh.file.Name.Name, cl); ok {
		if importPath, exists := fh.aliasImports[si.key.pkgAlias]; exists {
			key := UniqueIdentifierStructKey{
				ImportPath: importPath,
				StructName: si.key.name,
			}
			fh.structsInst[key] = append(fh.structsInst[key], si)
		}
	}
}

func (fh *FileHolder) CheckDotImports() {
	// check if there are some dotimports, this means struct insts with the same pkg
}

type UniqueIdentifierStructKey struct {
	ImportPath string
	StructName string
}
