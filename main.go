package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"main/internal"
	"path/filepath"
)

// enum(pending_as, runNing, completed, failed)
type IntState int

// enum(pending_as, runNing, completed, failed)
type StrState string

// enum(pending_as, runNing, completed, failed)
type PandState int

const UsageText = "Usage of enumer: \n" +
	"\t enumer -gen \n" +
	"For more information, see: \n" +
	"\t https://pkg.go.dev/golang.org/x/tools/cmd/stringer \n" +
	"Flags: \n"

var (
	generateFlag = flag.Bool("gen", false, "run generate")
	vendorFlag   = flag.Bool("vendor", false, "show detailed logs of execution")
)

func Usage() {
	fmt.Print(UsageText)
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("stringer: ")
	flag.Usage = Usage
	flag.Parse()

	if *generateFlag {
		DoGenerate()
		return
	}
	Usage()
}

func DoGenerate() {
	if *vendorFlag {
		log.Println("Begin scan packages...")
	}

	files, err := internal.GlobFiles()
	if err != nil {
		log.Fatal("glob error:", err)
	}
	fileSet := token.NewFileSet()
	var packageName string
	var allEnums []*internal.FutureEnum
	for _, filePath := range files {
		if *vendorFlag {
			log.Printf("Scan files: %s... \n", filePath)
		}
		absolutePath, err := filepath.Abs(filePath)
		if err != nil {
			if *vendorFlag {
				log.Printf("abs path err: %s \n", filePath)
			}
			continue
		}
		file, err := parser.ParseFile(fileSet, absolutePath, nil, parser.ParseComments)
		if err != nil {
			if *vendorFlag {
				log.Printf("parse file err: %s \n", absolutePath)
			}
			continue
		}
		if packageName == "" {
			packageName = file.Name.Name
		}
		enums, err := internal.ParseEnums(file)
		if err != nil {
			if *vendorFlag {
				log.Printf("parse file enums err: %s \n", filePath)
			}
			continue
		}
		allEnums = append(allEnums, enums...)
	}

	err = internal.GenerateEnumFile("./enumer.g.go", packageName, allEnums)
	if err != nil {
		log.Fatal(err)
	}
}
