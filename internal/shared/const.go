package shared

import (
	"fmt"
	"os"
	"path/filepath"
)

const ProjectLink = "https://github.com/oesand/go-enumer"

var KnownPackages = map[string]string{
	"fmt":   "\"fmt\"",
	"sql":   "\"database/sql\"",
	"cases": "\"github.com/oesand/go-enumer/cases\"",
	"sqlen": "\"github.com/oesand/go-enumer/sql\"",
	"types": "\"github.com/oesand/go-enumer/types\"",
}

type KnownEnumType string

const (
	IntEnum    KnownEnumType = "int"
	StringEnum KnownEnumType = "string"
)

var EnumSupportedTypes = map[string]KnownEnumType{
	"string": StringEnum,
	"int":    IntEnum,
	"int32":  IntEnum,
	"int64":  IntEnum,
}

func OpenFile(path string) (*os.File, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("fail to get absolute path fot output file: %s", err)
	}
	file, err := os.OpenFile(absolutePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
