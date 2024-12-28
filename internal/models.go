package internal

type File struct {
	Package string
	Enums   []*FutureEnum
}

type FutureEnum struct {
	TypeName string
	EnumName string
	Values   []EnumValue
}

type EnumValue struct {
	Name  string
	Value string
}
