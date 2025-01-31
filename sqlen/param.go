package sqlen

import (
	"database/sql"
	"fmt"
	"github.com/oesand/go-enumer/types"
)

var defaultFormatters = (func() map[string]ParamFormatter {
	formatters := types.PairSlice[ParamFormatter, []string]{
		{DollarFormatter, []string{"postgres", "pgx", "pq-timeouts", "cloudsqlpostgres", "ql", "nrpostgres", "cockroach"}},
		{QuestionFormatter, []string{"mysql", "sqlite3", "nrmysql", "nrsqlite3"}},
		{NamedFormatter, []string{"oci8", "ora", "goracle", "godror"}},
		{IndexedFormatter, []string{"sqlserver", "azuresql"}},
	}
	output := map[string]ParamFormatter{}
	for formatter, names := range formatters.I() {
		for _, name := range names {
			output[name] = formatter
		}
	}
	return output
})()

func DefaultFormatter(driverName string) ParamFormatter {
	if frm, has := defaultFormatters[driverName]; has {
		return frm
	}
	panic(fmt.Sprintf("unknown driver: %s", driverName))
}

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
	NamedFormatter    ParamFormatter = &paramFormatting{Prefix: "@param", FType: NamedFormatterType}
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

func namedParam(index int, value any) sql.NamedArg {
	return sql.Named(fmt.Sprintf("param%d", index), value)
}
