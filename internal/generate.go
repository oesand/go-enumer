package internal

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/oesand/go-enumer/cases"
	"github.com/oesand/go-enumer/internal/shared"
	"io"
	"strings"
	"text/template"
)

//go:embed template/*.tmpl
var content embed.FS

const fileHeaderText = "// Code generated by go-enumer[" + shared.ProjectLink + "]. DO NOT EDIT! \n"

func GenerateFile(filePath string, data *shared.GenerateData) error {
	if data.TotalCount() == 0 {
		return nil
	}
	file, err := shared.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\npackage %s\n", fileHeaderText, data.PackageName))
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"sumWithLen": func(one int, str string) int {
			return one + len(str)
		},
		"hasTag": func(tags map[string]string, tag string) bool {
			_, has := tags[tag]
			return has
		},
		"eqTag": func(tags map[string]string, tag string, value string) bool {
			tval, has := tags[tag]
			return has && tval == value
		},
	}

	if len(data.Imports) > 0 {
		var importsContent bytes.Buffer
		importsContent.WriteString("import (\n")
		importsMap := make(map[string]string, len(data.Imports))
		for i, imp := range data.Imports.Values() {
			alias := fmt.Sprintf("ial%d", i+1)
			importsMap[imp] = alias
			importsContent.WriteString(fmt.Sprintf("\t%s %s\n", alias, imp))
		}
		importsContent.WriteString(")\n")
		_, err = importsContent.WriteTo(file)
		if err != nil {
			return err
		}
		genType := func(info *shared.ExtraTypeInfo) string {
			var typeString strings.Builder
			if info.Starred {
				typeString.WriteRune('*')
			}
			if info.ImportPath != "" {
				typeString.WriteString(importsMap[info.ImportPath])
				typeString.WriteByte('.')
			}
			typeString.WriteString(info.TypeName)
			return typeString.String()
		}
		funcMap["genItemType"] = genType
		funcMap["genType"] = func(info *shared.ExtraTypeInfo) string {
			dcl := genType(info)
			if info.IsArray {
				dcl = fmt.Sprintf("[]%s", dcl)
			}
			return dcl
		}
		funcMap["knownAlias"] = func(name string) string {
			importPath, has := shared.KnownPackages[name]
			if !has {
				panic(fmt.Sprintf("unknown package alias: %s", name))
			}
			alias, has := importsMap[importPath]
			if !has {
				panic(fmt.Sprintf("unknown package import: %s", importPath))
			}
			return alias
		}
	} else {
		funcMap["knownAlias"] = func(name string) string {
			panic("no imports")
		}
	}

	if len(data.Enums) > 0 {
		err = executeTemplate(file, funcMap, "enum.tmpl", map[string]any{
			"Enums": data.Enums,
		})
		if err != nil {
			return err
		}
	}
	if len(data.Structs) > 0 {
		funcMap["genFieldNameAll"] = func(info *shared.StructInfo) string {
			var content strings.Builder
			for _, field := range info.Fields {
				content.WriteString(info.FieldCase.From(field.FieldName))
			}
			return content.String()
		}
		funcMap["camelCase"] = func(data string) string {
			return cases.ToCamelCase(data)
		}

		err = executeTemplate(file, funcMap, "struct.tmpl", map[string]any{
			"Structs": data.Structs,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func executeTemplate(w io.Writer, funcMap template.FuncMap, templateName string, data map[string]any) error {
	tmpl, err := template.New("").Funcs(funcMap).ParseFS(content, fmt.Sprintf("template/%s", templateName))
	if err != nil {
		return err
	}
	err = tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		return err
	}
	return nil
}
