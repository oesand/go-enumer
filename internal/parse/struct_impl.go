package parse

import (
	"fmt"
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
				case "query":
					definedTags[key] = ""
					knownImports.Add("cases")
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

	structInfo := &shared.StructInfo{
		Name:         name,
		KnownImports: knownImports,
		GenerateKind: generationKind,
		Tags:         definedTags,
	}

	return structInfo, nil
}
