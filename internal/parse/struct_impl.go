package parse

import (
	"fmt"
	"github.com/oesand/go-enumer/internal/shared"
	"regexp"
)

var structExp = regexp.MustCompile(`(?i)^\s*enumer:\s*(\S+)`)

func parseStructType(name string, comment string) (*shared.StructInfo, error) {
	matches := structExp.FindStringSubmatch(comment)
	if matches == nil {
		return nil, nil
	}

	structInfo := &shared.StructInfo{
		Name: name,
	}

	//appliedKeys := make(map[string]struct{})
	//enumEndIndex := enumExp.FindStringIndex(comment)[1]
	//sequencedText := strings.Trim(comment[enumEndIndex:], " \n")

	generationKind := shared.StructGenKind(matches[1])
	switch generationKind {
	case shared.BuilderGenKind:
		structInfo.RequireImports = true
	default:
		return nil, fmt.Errorf("unknown enumer generation kind: %s", generationKind)
	}
	structInfo.GenerateKind = generationKind

	return structInfo, nil
}
