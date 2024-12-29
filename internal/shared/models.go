package shared

type File struct {
	Package string
	Enums   []*EnumInfo
	Structs []*StructInfo
}

type GenerateData struct {
	Enums   []*EnumInfo
	Structs []*StructInfo
}

type EnumInfo struct {
	TypeName string
	EnumName string
	Values   []EnumValue
}

type EnumValue struct {
	Name  string
	Value string
}

type StructInfo struct{}
