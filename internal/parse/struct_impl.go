package parse

import (
	"fmt"
	"github.com/oesand/go-enumer/cases"
	"github.com/oesand/go-enumer/internal/shared"
	"github.com/oesand/go-enumer/types"
	"regexp"
	"strings"
)

var structExp = regexp.MustCompile(`(?i)^\s*enumer:\s*(\S+)`)

func parseStructType(name string, comment string) (*shared.StructInfo, error) {
	matches := structExp.FindStringSubmatch(comment)
	if matches == nil {
		return nil, nil
	}

	var fieldCase cases.CaseType
	var definedTags map[string]string
	var knownImports types.Set[string]
	generationKind := shared.StructGenKind(strings.ToLower(matches[1]))
	declEndIndex := structExp.FindStringIndex(comment)[1]
	switch generationKind {
	case shared.BuilderGenKind:
		definedTags = make(map[string]string)
		sequencedText := strings.Trim(comment[declEndIndex:], " \n")
		if sequencedText != "" {
			err := visitAllTags(sequencedText, false, func(key, value string) (err error) {
				if _, has := definedTags[key]; has {
					return fmt.Errorf("duplicated tag: %s", key)
				}
				switch key {
				case "repo":
					definedTags[key] = ""
					fieldCase = cases.CaseType(value)
					knownImports.Add("sql")
					knownImports.Add("sqlen")
				default:
					return fmt.Errorf("unknown tag name: %s", key)
				}
				return
			})
			if err != nil {
				return nil, err
			}
		}
		knownImports.Add("fmt")
	default:
		return nil, fmt.Errorf("unknown enumer generation kind: %s", generationKind)
	}

	if !fieldCase.IsValid() {
		return nil, fmt.Errorf("invalid field case: %s", fieldCase)
	}

	structInfo := &shared.StructInfo{
		Name:         name,
		FieldCase:    fieldCase,
		KnownImports: knownImports,
		GenerateKind: generationKind,
		Tags:         definedTags,
	}

	return structInfo, nil
}
