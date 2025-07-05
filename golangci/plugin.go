package structfieldinitorder

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/structfieldinitorder/analyzer"
)

func init() {
	register.Plugin("structfieldinitorder", New)
}

var _ register.LinterPlugin = new(structFieldInitOrderPlugin)

type structFieldInitOrderPlugin struct{}

func New(_ any) (register.LinterPlugin, error) {
	return &structFieldInitOrderPlugin{}, nil
}

func (p structFieldInitOrderPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.NewAnalyzer(),
	}, nil
}

func (p structFieldInitOrderPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
