package parse

import (
	"fmt"
	"github.com/oesand/go-enumer/internal/shared"
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
	requiredImports := make(map[string]struct{})
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
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			return true
		}
		for _, spec := range decl.Specs {
			tspec, ok := spec.(*ast.TypeSpec)
			if !ok {
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
					return true
				}
				if info == nil {
					continue
				}

				var fields []*shared.StructField
				for _, field := range typ.Fields.List {
					for _, fieldName := range field.Names {
						if !fieldName.IsExported() {
							continue
						}
						var starred bool
						var typeInfo *shared.ExtraTypeInfo
						fieldType := field.Type
					rollback:
						switch ftyp := fieldType.(type) {
						case *ast.Ident:
							typeInfo = &shared.ExtraTypeInfo{
								Starred:  starred,
								TypeName: ftyp.Name,
							}
						case *ast.SelectorExpr:
							typeInfo = &shared.ExtraTypeInfo{
								Starred:  starred,
								TypeName: ftyp.Sel.Name,
							}
							if info.RequireImports {
								importPath := importsMap[ftyp.X.(*ast.Ident).Name]
								requiredImports[importPath] = struct{}{}
								typeInfo.ImportPath = importPath
							}
						case *ast.StarExpr:
							starred = true
							fieldType = ftyp.X
							goto rollback
						}
						if typeInfo != nil {
							fields = append(fields, &shared.StructField{
								FieldName: fieldName.Name,
								TypeInfo:  typeInfo,
							})
						}
					}
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

				if _, has := shared.EnumSupportedTypes[typeName]; !has {
					return true
				}

				var enum *shared.EnumInfo
				enum, err = parseEnumType(typeName, name, doc)
				if err != nil {
					err = newLocatedErr(fileSet, filepath.Base(absolutePath), tspec, err.Error())
					return true
				}
				if enum == nil {
					continue
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
	imports := make([]string, 0, len(requiredImports))
	for importPath := range requiredImports {
		imports = append(imports, importPath)
	}
	fileInfo := &shared.ParsedFile{
		Package: file.Name.Name,
		Imports: imports,
		Items:   parsedItems,
	}
	return fileInfo, nil
}
