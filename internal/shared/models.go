package shared

type ItemType int

const (
	EnumItemType ItemType = iota
	StructItemType
)

type ParsedFile struct {
	Package string
	Imports []string
	Items   []*ParsedItem
}

type ParsedItem struct {
	ItemType ItemType
	Enum     *EnumInfo
	Struct   *StructInfo
}

type GenerateData struct {
	PackageName string
	Imports     []string
	Enums       []*EnumInfo
	Structs     []*StructInfo
}

func (g *GenerateData) TotalCount() int {
	return len(g.Enums) + len(g.Structs)
}

type EnumInfo struct {
	TypeName string
	EnumName string
	Values   []*EnumValue
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
	Name   string
	Fields []*StructField

	RequireImports bool
	GenerateKind   StructGenKind
}

type StructField struct {
	FieldName string
	TypeInfo  *ExtraTypeInfo
}

type ExtraTypeInfo struct {
	Starred    bool
	ImportPath string
	TypeName   string
}
