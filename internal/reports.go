package internal

import (
	"fmt"
	"slices"

	"golang.org/x/tools/go/analysis"
)

func ReportIfStructFieldsNotInOrder(pass *analysis.Pass, structSpec *StructSpec, structInit *StructInst) {
	instantiatedFieldNames := structInit.GetFieldNames()
	expectedFieldOrder := slices.DeleteFunc(structSpec.GetFieldNames(), func(s string) bool {
		return !slices.Contains(instantiatedFieldNames, s)
	})
	for i := range instantiatedFieldNames {
		if instantiatedFieldNames[i] != expectedFieldOrder[i] {
			pass.Report(analysis.Diagnostic{
				Pos: structInit.Pos(),
				End: structInit.End(),
				Message: fmt.Sprintf("fields for struct %q are not instantiated in order",
					structSpec.Name),
				URL: "https://github.com/manuelarte/structfieldinitorder?tab=readme-ov-file",
				// TODO(manuelarte): propose fix
			})
			return
		}
	}
}
