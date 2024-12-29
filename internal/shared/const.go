package shared

import (
	"fmt"
	"os"
	"path/filepath"
)

const ProjectLink = "https://github.com/oesand/go-enumer"

var EnumSupportedTypes = map[string]string{
	"string": "string",
	"int":    "int",
	"int32":  "int",
	"int64":  "int",
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
