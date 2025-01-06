package sqlen

import (
	"fmt"
	"github.com/oesand/go-enumer/types"
	"strings"
)

var defaultFormatters = (func() map[string]ParamFormatter {
	formatters := []types.Tuple[ParamFormatter, []string]{
		{DollarFormatter, []string{"postgres", "pgx", "pq-timeouts", "cloudsqlpostgres", "ql", "nrpostgres", "cockroach"}},
		{QuestionFormatter, []string{"mysql", "sqlite3", "nrmysql", "nrsqlite3"}},
		{NamedFormatter, []string{"oci8", "ora", "goracle", "godror"}},
		{IndexedFormatter, []string{"sqlserver", "azuresql"}},
	}
	output := map[string]ParamFormatter{}
	for _, data := range formatters {
		for _, name := range data.Second {
			output[name] = data.First
		}
	}
	return output
})()

type FormatterType int

const (
	SymbolFormatterType FormatterType = iota
	IndexedFormatterType
	NamedFormatterType
)

type ParamFormatter interface {
	Type() FormatterType
	Format(index int) string
}

var (
	QuestionFormatter ParamFormatter = &paramFormatting{Prefix: "?", FType: SymbolFormatterType}
	DollarFormatter   ParamFormatter = &paramFormatting{Prefix: "$", FType: IndexedFormatterType}
	IndexedFormatter  ParamFormatter = &paramFormatting{Prefix: "@p", FType: IndexedFormatterType}
	NamedFormatter    ParamFormatter = &paramFormatting{Prefix: "@prm", FType: NamedFormatterType}
)

type paramFormatting struct {
	Prefix string
	FType  FormatterType
}

func (p *paramFormatting) Type() FormatterType {
	return p.FType
}

func (p *paramFormatting) Format(index int) string {
	if p.FType == SymbolFormatterType {
		return p.Prefix
	}
	return fmt.Sprintf("%s%d", p.Prefix, index)
}

func repeatParamPlaceholders(format ParamFormatter, offset, count int) string {
	var valueString strings.Builder
	for i := 0; i < count; i++ {
		if i > 0 {
			valueString.WriteString(", ")
		}
		valueString.WriteString(format.Format(offset + i + 1))
	}
	return valueString.String()
}
