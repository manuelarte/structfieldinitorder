package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/manuelarte/structfieldinitorder/analyzer"
)

func main() {
	multichecker.Main(analyzer.NewAnalyzer())
}
