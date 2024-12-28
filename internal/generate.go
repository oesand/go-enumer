package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

var supportedTypes = map[string]string{
	"string": "string",
	"int":    "int",
	"int32":  "int",
	"int64":  "int",
}

func GenerateFiles(packageName string, enums []*FutureEnum) error {
	absolutePath, err := filepath.Abs("./enumer.g.go")
	if err != nil {
		return fmt.Errorf("fail to get absolute path fot output file: %s", err)
	}
	enumerFile, err := os.OpenFile(absolutePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer enumerFile.Close()
	err = generateEnumerFileContent(packageName, enums, enumerFile)
	return err
}
