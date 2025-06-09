package internal

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func reportStructFieldsNotInOrder(pass *analysis.Pass, structSpec *ast.TypeSpec) {
	// TODO(manuelarte)
}

