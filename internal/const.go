package internal

import (
	"fmt"
	"go/ast"
	"go/token"
)

const ProjectLink = "https://github.com/oesand/go-enumer"

func newLocatedErr(fileSet *token.FileSet, fileName string, tspec *ast.TypeSpec, text string, a ...any) error {
	pos := tspec.Pos()
	file := fileSet.File(pos)
	line := file.Line(pos)
	return &locatedError{
		fileName: fileName,
		typeName: tspec.Name.Name,
		line:     line,
		text:     fmt.Sprintf(text, a...),
	}
}

type locatedError struct {
	fileName string
	typeName string
	line     int
	text     string
}

func (e *locatedError) Error() string {
	return fmt.Sprintf("[%s:%d @%s]: %s", e.fileName, e.line, e.typeName, e.text)
}

func allValuesUnique[S ~[]E, E comparable](slice S) bool {
	seen := make(map[E]struct{})
	for _, val := range slice {
		if _, has := seen[val]; has {
			return false
		}
		seen[val] = struct{}{}
	}
	return true
}
