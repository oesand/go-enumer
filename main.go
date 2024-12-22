package main

import (
	"flag"
	"fmt"
	"github.com/oesand/go-enumer/internal"
	"go/token"
	"log"
	"path/filepath"
	"strings"
)

const UsageText = "Usage of enumer: \n" +
	"\t go-enumer # Help - you here ;) \n" +
	"\t go-enumer gen # Generates enums from files current directory \n" +
	"For more information, see: \n" +
	"\t " + internal.ProjectLink + " \n"

func PrintUsage() {
	fmt.Print(UsageText)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("go-enumer: ")
	flag.Usage = PrintUsage
	flag.Parse()

	if flag.Arg(0) == "gen" {
		DoGenerate()
		return
	}
	PrintUsage()
}

func DoGenerate() {
	files, err := internal.GlobFiles()
	if err != nil {
		log.Fatal("glob error:", err)
	}
	fileSet := token.NewFileSet()
	var packageName string
	var allEnums []*internal.FutureEnum
	for _, fileName := range files {
		if strings.Count(fileName, ".") > 1 {
			continue
		}
		absolutePath, err := filepath.Abs(fileName)
		if err != nil {
			continue
		}
		file, err := internal.ParseFile(fileSet, absolutePath)
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
		allEnums = append(allEnums, file.Enums...)
	}

	log.Printf("generate file enumer.g.go with %d enums total", len(allEnums))
	err = internal.GenerateEnumFile("./enumer.g.go", packageName, allEnums)
	if err != nil {
		log.Fatal(err)
	}
}
