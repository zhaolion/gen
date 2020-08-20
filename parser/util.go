package parser

import (
	"go/ast"
	tc "go/types"
	"regexp"
	"strings"

	"github.com/zhaolion/gen/types"
)

// This clarifies when a pkg path has been canonicalized.
type importPathString string

func (s importPathString) String() string {
	return string(s)
}

// parsedFile is for tracking files with name
type parsedFile struct {
	name string
	file *ast.File
}

// key type for finding comments.
type fileLine struct {
	file string
	line int
}

type importAdapter struct {
	b *Builder
}

func (a importAdapter) Import(path string) (*tc.Package, error) {
	return a.b.importPackage(path, false)
}

// canonicalizeImportPath takes an import path and returns the actual package.
// It doesn't support nested vendoring.
func canonicalizeImportPath(importPath string) importPathString {
	if !strings.Contains(importPath, "/vendor/") {
		return importPathString(importPath)
	}

	return importPathString(importPath[strings.Index(importPath, "/vendor/")+len("/vendor/"):])
}

// regexErrPackageNotFound helps test the expected error for not finding a package.
var regexErrPackageNotFound = regexp.MustCompile(`^unable to import ".*?":.*`)

func isErrPackageNotFound(err error) bool {
	return regexErrPackageNotFound.MatchString(err.Error())
}

func splitLines(str string) []string {
	return strings.Split(strings.TrimRight(str, "\n"), "\n")
}

func tcNameToName(in string) types.Name {
	// Detect anonymous type names. (These may have '.' characters because
	// embedded types may have packages, so we detect them specially.)
	if strings.HasPrefix(in, "struct{") ||
		strings.HasPrefix(in, "<-chan") ||
		strings.HasPrefix(in, "chan<-") ||
		strings.HasPrefix(in, "chan ") ||
		strings.HasPrefix(in, "func(") ||
		strings.HasPrefix(in, "*") ||
		strings.HasPrefix(in, "map[") ||
		strings.HasPrefix(in, "[") {
		return types.Name{Name: in}
	}

	// Otherwise, if there are '.' characters present, the name has a
	// package path in front.
	nameParts := strings.Split(in, ".")
	name := types.Name{Name: in}
	if n := len(nameParts); n >= 2 {
		// The final "." is the name of the type--previous ones must
		// have been in the package path.
		name.Package, name.Name = strings.Join(nameParts[:n-1], "."), nameParts[n-1]
	}
	return name
}

func tcFuncNameToName(in string) types.Name {
	name := strings.TrimPrefix(in, "func ")
	nameParts := strings.Split(name, "(")
	return tcNameToName(nameParts[0])
}

func tcVarNameToName(in string) types.Name {
	nameParts := strings.Split(in, " ")
	// nameParts[0] is "var".
	// nameParts[2:] is the type of the variable, we ignore it for now.
	return tcNameToName(nameParts[1])
}
