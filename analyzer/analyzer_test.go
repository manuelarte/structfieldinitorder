package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testCases := []struct {
		desc     string
		patterns string
	}{
		{
			desc:     "default",
			patterns: "simple/...",
		},
		{
			desc:     "imports",
			patterns: "imports/...",
		},
		{
			desc:     "dot imports",
			patterns: "dotimports/...",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			a := NewAnalyzer()

			analysistest.Run(t, analysistest.TestData(), a, test.patterns)
		})
	}
}
