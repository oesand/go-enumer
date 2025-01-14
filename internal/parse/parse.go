package parse

import (
	"errors"
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
			case *ast.StructType:
				name := tspec.Name.Name
				doc := decl.Doc.Text()
				var info *shared.StructInfo
				info, err = parseStructType(name, doc)
				if err != nil {
					err = newLocatedErr(fileSet, filepath.Base(absolutePath), tspec, err.Error())
					return false
				}
				if info == nil {
					continue
				}

				for importName := range info.KnownImports {
					importPath, has := shared.KnownPackages[importName]
					if !has {
						panic(fmt.Sprintf("unknown package predefined alias: %s", importName))
					}
					requiredImports.Add(importPath)
				}

				var fields []*shared.StructField
				for _, field := range typ.Fields.List {
					for _, fieldName := range field.Names {
						if !fieldName.IsExported() {
							continue
						}
						fieldType := field.Type
						var typeInfo shared.ExtraTypeInfo
					rollback:
						switch ftyp := fieldType.(type) {
						case *ast.Ident:
							typeInfo.TypeName = ftyp.Name
						case *ast.SelectorExpr:
							typeInfo.TypeName = ftyp.Sel.Name
							importPath := importsMap[ftyp.X.(*ast.Ident).Name]
							requiredImports.Add(importPath)
							typeInfo.ImportPath = importPath
						case *ast.StarExpr:
							typeInfo.Starred = true
							fieldType = ftyp.X
							goto rollback
						case *ast.ArrayType:
							typeInfo.IsArray = true
							fieldType = ftyp.Elt
							goto rollback
						default:
							err = fmt.Errorf("unsupported field type: %T\n", fieldType)
						}
						if err != nil {
							break
						}
						fields = append(fields, &shared.StructField{
							FieldName: fieldName.Name,
							CasedName: info.FieldCase.From(fieldName.Name),
							TypeInfo:  &typeInfo,
						})
					}
					if err != nil {
						break
					}
				}
				if len(fields) == 0 {
					err = errors.New("cannot generate builder, empty fields")
					return false
				}
				if err != nil {
					err = newLocatedErr(fileSet, filepath.Base(absolutePath), tspec, err.Error())
					return false
				}
				info.Fields = fields
				parsedItems = append(parsedItems, &shared.ParsedItem{
					ItemType: shared.StructItemType,
					Struct:   info,
				})

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
				requiredImports.Add(shared.KnownPackages["fmt"])

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
