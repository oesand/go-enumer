package internal

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	enumExp = regexp.MustCompile(`(?i)^\s*enum\(([^)]*)\)`)
)

func GlobFiles() ([]string, error) {
	matches, err := filepath.Glob("./*.go")
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func ParseEnums(file *ast.File) ([]*FutureEnum, error) {
	var enums []*FutureEnum
	ast.Inspect(file, func(node ast.Node) bool {
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			return true
		}
		for _, spec := range decl.Specs {
			tspec := spec.(*ast.TypeSpec)

			tp := tspec.Type.(*ast.Ident).Name
			name := tspec.Name.Name
			doc := decl.Doc.Text()
			enum, _ := parseType(tp, name, doc)
			enums = append(enums, enum)
		}
		return false
	})
	return enums, nil
}

func parseType(typeName string, name string, comment string) (*FutureEnum, error) {
	if _, has := supportedTypes[typeName]; !has {
		return nil, fmt.Errorf("not supported type(%s)", typeName)
	}
	matches := enumExp.FindStringSubmatch(comment)
	if matches == nil {
		return nil, nil
	}
	valuesString := strings.ReplaceAll(matches[1], " ", "")
	if valuesString == "" {
		return nil, errors.New("empty enum values, see examples")
	}
	values := strings.Split(valuesString, ",")

	enumInfo := &FutureEnum{
		TypeName:   supportedTypes[typeName],
		EnumName:   name,
		ValueNames: values,
	}
	return enumInfo, nil
}
