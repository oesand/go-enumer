package internal

import (
	"go/ast"
)

type File struct {
	Name string
	file *ast.File // Parsed AST.
}

type FutureEnum struct {
	// Main info
	TypeName   string
	EnumName   string
	ValueNames []string

	// Extra options
	reversedName bool
}
