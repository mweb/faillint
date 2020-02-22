package faillint_test

import (
	"testing"

	"github.com/fatih/faillint/faillint"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()

	var tests = []struct {
		name        string
		paths       string
		testpaths   string
	}{
		{
			name:        "a",
			paths:       "errors",
			testpaths:   "",
		},
		{
			name:        "b",
			paths:       "",
			testpaths:   "",
		},
		{
			name:        "c",
			paths:       "errors=", // malformed suggestion
			testpaths:   "",
		},
		{
			name:        "d",
			paths:       "errors=github.com/pkg/errors",
			testpaths:   "",
		},
		{
			name:        "e",
			paths:       "errors=github.com/pkg/errors,golang.org/x/net/context=context",
			testpaths:   "",
		},
		{
			name:        "f",
			paths:       "errors",
			testpaths:   "log",
		},
		{
			name:        "g",
			paths:       "errors",
			testpaths:   "log",
		},
		{
			name:        "h",
			paths:       "errors",
			testpaths:   "",
		},
		{
			name:        "i",
			paths:       "reflect,errors",
			testpaths:   "errors",
		},
	}
	for _, ts := range tests {
		ts := ts
		t.Run(ts.name, func(t *testing.T) {
			a := faillint.NewAnalyzer()
			a.Flags.Set("paths", ts.paths)
			a.Flags.Set("testpaths", ts.testpaths)
			analysistest.Run(t, testdata, a, ts.name)
		})
	}
}
