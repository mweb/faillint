// Package faillint defines an Analyzer that fails when a package is imported
// that matches against a set of defined packages.
package faillint

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type faillint struct {
	paths     string // -paths flag
	testpaths string // -testpaths flag
}

// NewAnalyzer create a faillint analyzer
func NewAnalyzer() *analysis.Analyzer {
	f := faillint{
		paths:     "",
		testpaths: "",
	}
	a := &analysis.Analyzer{
		Name:             "faillint",
		Doc:              "report unwanted import path usages",
		Run:              f.run,
		RunDespiteErrors: true,
	}
	a.Flags.StringVar(&f.paths, "paths", "", "import paths to fail")
	a.Flags.StringVar(&f.testpaths, "testpaths", "", "import paths for all _test.go files and packages to fail")
	return a
}

// Run is the runner for an analysis pass
func (f *faillint) run(pass *analysis.Pass) (interface{}, error) {
	if f.paths == "" {
		return nil, nil
	}

	imports, suggestions := getRules(strings.Split(f.paths, ","))
	testImports, testSuggestions := imports, suggestions
	if f.testpaths != "" {
		testImports, testSuggestions = getRules(strings.Split(f.testpaths, ","))
	}

	for _, file := range pass.Files {
		if strings.Contains(pass.Fset.File(file.Package).Name(), "_test.go") {
			checkImport(pass, file, testImports, testSuggestions)
		} else {
			checkImport(pass, file, imports, suggestions)
		}

	}

	return nil, nil
}

// checkImport checks the import statement of a file against a list with
// forbidden imports
func checkImport(pass *analysis.Pass, file *ast.File, imports []string, suggestions map[string]string) {
	for _, path := range imports {
		imp := usesImport(file, path)
		if imp == nil {
			continue
		}

		impPath := importPath(imp)

		msg := fmt.Sprintf("package %q shouldn't be imported", impPath)
		if s := suggestions[impPath]; s != "" {
			msg += fmt.Sprintf(", suggested: %q", s)
		}

		pass.Reportf(imp.Path.Pos(), msg)
	}
}

// getRules generates a list with all the imports and a map with all the
// suggestions from the paths strings
func getRules(paths []string) ([]string, map[string]string) {
	suggestions := make(map[string]string, len(paths))
	imports := make([]string, 0, len(paths))

	for _, s := range paths {
		imps := strings.Split(s, "=")

		imp := imps[0]
		suggest := ""
		if len(imps) == 2 {
			suggest = imps[1]
		}

		imports = append(imports, imp)
		suggestions[imp] = suggest
	}

	return imports, suggestions
}

// usesImport reports whether a given import is used.
func usesImport(f *ast.File, path string) *ast.ImportSpec {
	spec := importSpec(f, path)
	if spec == nil {
		return nil
	}

	name := spec.Name.String()
	switch name {
	case "<nil>":
		// If the package name is not explicitly specified,
		// make an educated guess. This is not guaranteed to be correct.
		lastSlash := strings.LastIndex(path, "/")
		if lastSlash == -1 {
			name = path
		} else {
			name = path[lastSlash+1:]
		}
	case "_", ".":
		// Not sure if this import is used - err on the side of caution.
		return spec
	}

	var used bool
	ast.Inspect(f, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if ok && isTopName(sel.X, name) {
			used = true
		}
		return true
	})

	if used {
		return spec
	}

	return nil
}

// importSpec returns the import spec if f imports path,
// or nil otherwise.
func importSpec(f *ast.File, path string) *ast.ImportSpec {
	for _, s := range f.Imports {
		if importPath(s) == path {
			return s
		}
	}
	return nil
}

// importPath returns the unquoted import path of s,
// or "" if the path is not properly quoted.
func importPath(s *ast.ImportSpec) string {
	t, err := strconv.Unquote(s.Path.Value)
	if err == nil {
		return t
	}
	return ""
}

// isTopName returns true if n is a top-level unresolved identifier with the given name.
func isTopName(n ast.Expr, name string) bool {
	id, ok := n.(*ast.Ident)
	return ok && id.Name == name && id.Obj == nil
}
