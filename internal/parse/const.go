package parse

import (
	"fmt"
	"go/ast"
	"go/token"
)

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
