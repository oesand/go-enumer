package internal

type File struct {
	Package string
	Enums   []*FutureEnum
}

type FutureEnum struct {
	// Main info
	TypeName   string
	EnumName   string
	ValueNames []string

	// Extra options
	reversedName bool
	prefix       string
}
