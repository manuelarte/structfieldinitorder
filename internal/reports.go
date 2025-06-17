package internal

import (
	"fmt"
	"slices"

	"golang.org/x/tools/go/analysis"
)

func ReportIfStructFieldsNotInOrder(pass *analysis.Pass, structSpecs *StructSpecs, si IStructInst) {
	instantiatedFieldNames := si.GetFieldNames()
	expectedFieldOrder := slices.DeleteFunc(structSpecs.GetFieldNames(), func(s string) bool {
		return !slices.Contains(instantiatedFieldNames, s)
	})
	for i := range instantiatedFieldNames {
		if instantiatedFieldNames[i] != expectedFieldOrder[i] {
			pass.Report(analysis.Diagnostic{
				Pos: si.Pos(),
				End: si.End(),
				Message: fmt.Sprintf("fields for struct %q are not instantiated in order",
					structSpecs.Name),
				URL: "https://github.com/manuelarte/structfieldinitorder?tab=readme-ov-file",
				// TODO(manuelarte): propose fix
			})
			return
		}
	}
}
