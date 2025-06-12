package internal

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func reportStructFieldsNotInOrder(pass *analysis.Pass, structSpec *ast.TypeSpec, structInit *StructInit) {
	pass.Report(analysis.Diagnostic{
		Pos: structInit.Lbrace,
		Message: fmt.Sprintf("fields for struct %q are not instantiated in order",
			structSpec.Name),
		URL: "https://github.com/manuelarte/structfieldinitorder?tab=readme-ov-file",
		// TODO(manuelarte): propose fix
	})
}
