package shared

import (
	"github.com/oesand/go-enumer/cases"
	"github.com/oesand/go-enumer/types"
)

type ItemType int

const (
	EnumItemType ItemType = iota
	StructItemType
)

type ParsedFile struct {
	Package string
	Imports types.Set[string]
	Items   []*ParsedItem
}

type ParsedItem struct {
	ItemType ItemType
	Enum     *EnumInfo
	Struct   *StructInfo
}

type GenerateData struct {
	PackageName string
	Imports     types.Set[string]
	Enums       []*EnumInfo
	Structs     []*StructInfo
}

func (g *GenerateData) TotalCount() int {
	return len(g.Enums) + len(g.Structs)
}

type EnumInfo struct {
	TypeName KnownEnumType
	EnumName string
	Values   []*EnumValue
	Tags     map[string]string
}

type EnumValue struct {
	Name  string
	Value string
}

type StructGenKind string

const (
	BuilderGenKind StructGenKind = "builder"
)

type StructInfo struct {
	Name      string
	FieldCase cases.CaseType
	Fields    []*StructField

	KnownImports types.Set[string]
	GenerateKind StructGenKind
	Tags         map[string]string
}

type StructField struct {
	FieldName string
	CasedName string
	TypeInfo  *ExtraTypeInfo
}

type ExtraTypeInfo struct {
	IsArray    bool
	Starred    bool
	ImportPath string
	TypeName   string
}
