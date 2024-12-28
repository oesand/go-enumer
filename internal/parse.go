package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
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

func ParseFile(fileSet *token.FileSet, absolutePath string) (*File, error) {
	file, err := parser.ParseFile(fileSet, absolutePath, nil, parser.ParseComments)
	if err != nil {
		err = fmt.Errorf("cannot parse file \"%s\": %v", filepath.Base(absolutePath), err)
		return nil, err
	}

	var enums []*FutureEnum
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
			var enum *FutureEnum
			enum, err = parseType(typeName, name, doc)
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
	fileInfo := &File{
		Package: file.Name.Name,
		Enums:   enums,
	}
	return fileInfo, nil
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
		return nil, fmt.Errorf("empty enum values, see examples %s", ProjectLink)
	}
	valueNames := strings.Split(valuesString, ",")

	var inverseNameOption bool
	prefixOption := name
	appliedKeys := make(map[string]struct{})

	enumEndIndex := enumExp.FindStringIndex(comment)[1]
	sequencedText := strings.Trim(comment[enumEndIndex:], " \n")
	if sequencedText != "" {
		matches := tagsExp.FindAllStringSubmatch(sequencedText, -1)
		if matches != nil {
			for _, match := range matches {
				var key string
				if match[1] == "" && match[2] == "" {
					key = match[0]
					switch key {
					case "inverse":
						inverseNameOption = true
					default:
						return nil, fmt.Errorf("unknown tag name: %s", match[0])
					}
				} else {
					key = match[1]
					value := match[2]
					switch key {
					case "prefix":
						prefixOption = toPascalCase(value)
					default:
						return nil, fmt.Errorf("unknown tag name: %s", key)
					}
				}
				if _, has := appliedKeys[key]; has {
					return nil, fmt.Errorf("duplicated tag: %s", key)
				}
				appliedKeys[key] = struct{}{}
			}
		}
	}

	values := make([]EnumValue, len(valueNames))

	for i, value := range valueNames {
		var name string
		if inverseNameOption {
			name = toPascalCase(value) + prefixOption
			if !unicode.IsLetter(rune(value[0])) {
				return nil, fmt.Errorf("generated invalid name for enum value(%s) with 'inverse' tag", value)
			}
		} else {
			name = prefixOption + toPascalCase(value)
		}

		values[i] = EnumValue{
			Name:  name,
			Value: value,
		}
	}

	enumInfo := &FutureEnum{
		TypeName: supportedTypes[typeName],
		EnumName: name,
		Values:   values,
	}

	return enumInfo, nil
}
