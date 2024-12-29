package main

import (
	"flag"
	"fmt"
	"github.com/oesand/go-enumer/internal/parse"
	"github.com/oesand/go-enumer/internal/shared"
	"go/token"
	"log"
	"path/filepath"
	"strings"
)

// enum(pending, running, completed)
type IntStatus int

// enum(pending, running, completed)
type StrStatus string

// @cls
type DataStr struct {
	Name string
}

const UsageText = "Usage of enumer: \n" +
	"\t go-enumer # Help - you here ;) \n" +
	"\t go-enumer gen # Generates enums from files current directory \n" +
	"For more information, see: \n" +
	"\t " + shared.ProjectLink + " \n"

func PrintUsage() {
	fmt.Print(UsageText)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("go-enumer: ")
	flag.Usage = PrintUsage
	flag.Parse()

	if flag.Arg(0) == "gen" {
		doGenerate()
		return
	}
	PrintUsage()
}

func doGenerate() {
	files, err := parse.GlobFiles()
	if err != nil {
		log.Fatal("glob error:", err)
	}
	fileSet := token.NewFileSet()
	var packageName string
	var generateData shared.GenerateData
	for _, fileName := range files {
		if strings.Count(fileName, ".") > 1 {
			continue
		}
		absolutePath, err := filepath.Abs(fileName)
		if err != nil {
			continue
		}
		file, err := parse.ParseFile(fileSet, absolutePath)
		if err != nil {
			log.Fatal(err)
		}
		if len(file.Enums) == 0 {
			continue
		}
		var enumsString strings.Builder
		for i, enum := range file.Enums {
			if i > 0 {
				enumsString.WriteString(", ")
			}
			enumsString.WriteString(enum.EnumName)
		}
		log.Printf("parsed file: %s [%s]", fileName, enumsString.String())
		if packageName == "" {
			packageName = file.Package
		}
		generateData.Enums = append(generateData.Enums, file.Enums...)
	}
	//if len(generateData.Enums) == 0 {
	//	log.Printf("file generation skipped, no enums found")
	//	return
	//}
	log.Printf("generate file enumer.g.go with %d enums total", len(generateData.Enums))
	//err = internal.GenerateEnumerFile(packageName, allEnums)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
