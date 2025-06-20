package internal

import (
	"bytes"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"go/ast"
	"go/format"
	"go/token"
	"maps"
	"os"
	"slices"

	"golang.org/x/tools/go/analysis"
)

func ReportIfStructFieldsNotInOrder(pass *analysis.Pass, structSpecs *StructSpecs, si StructInst) {

	// Create a new decorator, which will track the mapping between ast and dst nodes
	dec := decorator.NewDecorator(pass.Fset)

	kvs := si.GetKeyValueExpr()
	if slices.ContainsFunc(kvs, func(kv *ast.KeyValueExpr) bool {
		_, ok := kv.Key.(*ast.Ident)
		return !ok
	}) {
		// struct insts does not contains keys that are ast.Ident
		return
	}
	instFieldsIndexByFieldName := indexKeyValueByIdentName(kvs)
	instFieldNames := mapKeys(instFieldsIndexByFieldName)
	expectedFieldOrder := slices.DeleteFunc(structSpecs.GetFieldNames(), func(s string) bool {
		return !slices.Contains(instFieldNames, s)
	})
	var report bool
	sortedKeyValueOrder := make([]*ast.KeyValueExpr, len(kvs))
	for i := range instFieldNames {
		if instFieldNames[i] != expectedFieldOrder[i] {
			report = true
		}
		sortedKeyValueOrder[i] = instFieldsIndexByFieldName[expectedFieldOrder[i]]
	}
	if report {
		cl := si.GetCompositeLit()
		dstCl, _ := dec.DecorateNode(cl)
		dst.Fprint(os.Stdout, dstCl, ast.NotNilFilter)
		diag := analysis.Diagnostic{
			Pos: cl.Pos(),
			End: cl.End(),
			Message: fmt.Sprintf("fields for struct %q are not instantiated in order",
				structSpecs.Name),
			URL: "https://github.com/manuelarte/structfieldinitorder?tab=readme-ov-file",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message:   "Sorted struct instantiation fields.",
					TextEdits: toTextEdits(pass.Fset, sortedKeyValueOrder, cl, dstCl.(*dst.CompositeLit)),
				},
			},
		}
		pass.Report(diag)
	}
}

func indexKeyValueByIdentName(kvs []*ast.KeyValueExpr) map[string]*ast.KeyValueExpr {
	toReturn := make(map[string]*ast.KeyValueExpr, len(kvs))
	for _, kv := range kvs {
		//nolint:errcheck // impossible case
		ident := kv.Key.(*ast.Ident)
		toReturn[ident.Name] = kv
	}
	return toReturn
}

func mapKeys(myMap map[string]*ast.KeyValueExpr) []string {
	keys := make([]string, len(myMap))
	var i int
	for k := range maps.Keys(myMap) {
		keys[i] = k
		i++
	}
	return keys
}

func toTextEdits(fset *token.FileSet, sorted []*ast.KeyValueExpr, original *ast.CompositeLit, originalWithComments *dst.CompositeLit) []analysis.TextEdit {
	toReturn := make([]analysis.TextEdit, len(sorted))
	for i, kv := range sorted {
		var buf bytes.Buffer
		dstOrigin := originalWithComments.Elts[i].(*dst.KeyValueExpr)
		dst.Fprint(os.Stdout, dstOrigin, ast.NotNilFilter)
		origin := original.Elts[i].(*ast.KeyValueExpr)
		// TODO(manuelarte): handle error later
		_ = format.Node(&buf, fset, kv)
		toReturn[i] = analysis.TextEdit{
			Pos:     origin.Pos(),
			End:     origin.End(),
			NewText: buf.Bytes(),
		}
	}
	return toReturn
}
