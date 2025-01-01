package parse

import (
	"fmt"
	"github.com/oesand/go-enumer/internal/shared"
	"regexp"
	"strings"
)

var structExp = regexp.MustCompile(`(?i)^\s*enumer:\s*(\S+)`)

func parseStructType(name string, comment string) (*shared.StructInfo, error) {
	matches := structExp.FindStringSubmatch(comment)
	if matches == nil {
		return nil, nil
	}

	var requireImports bool
	var definedTags map[string]string
	generationKind := shared.StructGenKind(strings.ToLower(matches[1]))
	declEndIndex := structExp.FindStringIndex(comment)[1]
	switch generationKind {
	case shared.BuilderGenKind:
		requireImports = true
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
				default:
					return fmt.Errorf("unknown tag name: %s", key)
				}
				return
			})
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("unknown enumer generation kind: %s", generationKind)
	}

	structInfo := &shared.StructInfo{
		Name:           name,
		RequireImports: requireImports,
		GenerateKind:   generationKind,
		Tags:           definedTags,
	}

	return structInfo, nil
}
