package shared

import (
	"github.com/oesand/go-enumer/types"
)

type ItemType int

const (
	EnumItemType ItemType = iota
)

type ParsedFile struct {
	Package string
	Imports types.Set[string]
	Items   []*ParsedItem
}

type ParsedItem struct {
	ItemType ItemType
	Enum     *EnumInfo
}

type GenerateData struct {
	PackageName string
	Imports     types.Set[string]
	Enums       []*EnumInfo
}

func (g *GenerateData) TotalCount() int {
	return len(g.Enums)
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
