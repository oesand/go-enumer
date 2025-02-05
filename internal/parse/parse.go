package parse

import (
	"fmt"
	"github.com/oesand/go-enumer/internal/shared"
	"github.com/oesand/go-enumer/types"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

func GlobFiles() ([]string, error) {
	matches, err := filepath.Glob("*.go")
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func ParseFile(fileSet *token.FileSet, absolutePath string) (*shared.ParsedFile, error) {
	file, err := parser.ParseFile(fileSet, absolutePath, nil, parser.ParseComments)
	if err != nil {
		err = fmt.Errorf("cannot parse file \"%s\": %v", filepath.Base(absolutePath), err)
		return nil, err
	}
	var requiredImports types.Set[string]
	importsMap := make(map[string]string, len(file.Imports))
	for _, imp := range file.Imports {
		path := imp.Path.Value
		var alias string
		if imp.Name != nil {
			alias = imp.Name.Name
		} else {
			alias = filepath.Base(path[1 : len(path)-1])
		}
		importsMap[alias] = path
	}
	var parsedItems []*shared.ParsedItem
	ast.Inspect(file, func(node ast.Node) bool {
		if err != nil {
			return false
		}
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			return true
		}
		for _, spec := range decl.Specs {
			tspec, ok := spec.(*ast.TypeSpec)
			if !ok || !tspec.Name.IsExported() {
				continue
			}
			switch typ := tspec.Type.(type) {
			case *ast.Ident:
				typeName := typ.Name
				name := tspec.Name.Name
				doc := decl.Doc.Text()

				enumType, has := shared.EnumSupportedTypes[typeName]
				if !has {
					return false
				}

				var enum *shared.EnumInfo
				enum, err = parseEnumType(enumType, name, doc)
				if err != nil {
					err = newLocatedErr(fileSet, filepath.Base(absolutePath), tspec, err.Error())
					return false
				}
				if enum == nil {
					continue
				}
				if enum.TypeName == shared.IntEnum {
					requiredImports.Add(shared.KnownPackages["fmt"])
				}

				parsedItems = append(parsedItems, &shared.ParsedItem{
					ItemType: shared.EnumItemType,
					Enum:     enum,
				})
			}
		}
		return false
	})
	if err != nil {
		return nil, err
	}
	fileInfo := &shared.ParsedFile{
		Package: file.Name.Name,
		Imports: requiredImports,
		Items:   parsedItems,
	}
	return fileInfo, nil
}
