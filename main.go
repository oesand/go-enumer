package main

import (
	"flag"
	"fmt"
	"github.com/oesand/go-enumer/internal"
	"github.com/oesand/go-enumer/internal/parse"
	"github.com/oesand/go-enumer/internal/shared"
	"go/token"
	"log"
	"path/filepath"
	"strings"
)

type Jon struct{}

// enum(pending, running, completed)
type IntStatus int

// enum(pending, running, completed)
type StrStatus string

// enumer:builder
type DataStr struct {
	Name string
	Item *shared.ParsedItem
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
		if len(file.Items) == 0 {
			continue
		}
		if generateData.PackageName == "" {
			generateData.PackageName = file.Package
		}
		var parsedInfoString strings.Builder
		for i, item := range file.Items {
			if i > 0 {
				parsedInfoString.WriteString(", ")
			}
			switch item.ItemType {
			case shared.EnumItemType:
				parsedInfoString.WriteString(fmt.Sprintf("%s@enum", item.Enum.EnumName))
				generateData.Enums = append(generateData.Enums, item.Enum)
			case shared.StructItemType:
				parsedInfoString.WriteString(fmt.Sprintf("%s@struct", item.Struct.Name))
				generateData.Structs = append(generateData.Structs, item.Struct)
			default:
				log.Fatal("unknown item type:", item.ItemType)
			}
		}
		log.Printf("parsed file: %s [%s]", fileName, parsedInfoString.String())
		generateData.Imports = append(generateData.Imports, file.Imports...)
	}
	totalCount := generateData.TotalCount()
	if totalCount == 0 {
		log.Printf("file generation skipped, no enums found")
		return
	}
	log.Printf("generate file enumer.g.go with %d items total", totalCount)
	err = internal.GenerateFile(&generateData)
	if err != nil {
		log.Fatal(err)
	}
}
