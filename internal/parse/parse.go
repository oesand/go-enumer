package parse

import (
	"fmt"
	"github.com/oesand/go-enumer/internal/shared"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
)

var (
	enumExp = regexp.MustCompile(`(?i)^\s*enum\(([^)]*)\)`)
	tagsExp = regexp.MustCompile(`\b(\w+)\s*:\s*(\w+)\b|\b\w+\b`)
)

func GlobFiles() ([]string, error) {
	matches, err := filepath.Glob("*.go")
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func ParseFile(fileSet *token.FileSet, absolutePath string) (*shared.File, error) {
	file, err := parser.ParseFile(fileSet, absolutePath, nil, parser.ParseComments)
	if err != nil {
		err = fmt.Errorf("cannot parse file \"%s\": %v", filepath.Base(absolutePath), err)
		return nil, err
	}

	var enums []*shared.EnumInfo
	ast.Inspect(file, func(node ast.Node) bool {
		if err != nil {
			return true
		}
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			return true
		}
		for _, spec := range decl.Specs {
			tspec, sucs := spec.(*ast.TypeSpec)
			if !sucs {
				continue
			}
			ident, sucs := tspec.Type.(*ast.Ident)
			if !sucs {
				continue
			}

			typeName := ident.Name
			name := tspec.Name.Name
			doc := decl.Doc.Text()
			var enum *shared.EnumInfo
			enum, err = parseEnumType(typeName, name, doc)
			if err != nil {
				err = newLocatedErr(fileSet, filepath.Base(absolutePath), tspec, err.Error())
				return true
			}
			if enum == nil {
				continue
			}
			enums = append(enums, enum)
		}
		return false
	})
	if err != nil {
		return nil, err
	}
	fileInfo := &shared.File{
		Package: file.Name.Name,
		Enums:   enums,
	}
	return fileInfo, nil
}
