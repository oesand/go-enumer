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
	//if sequencedText != "" {
	//	matches := tagsExp.FindAllStringSubmatch(sequencedText, -1)
	//	if matches != nil {
	//		for _, match := range matches {
	//			var key string
	//			if match[1] == "" && match[2] == "" {
	//				key = match[0]
	//				switch key {
	//				case "builder":
	//					structInfo.RequireImports = true
	//				default:
	//					return nil, fmt.Errorf("unknown tag name: %s", match[0])
	//				}
	//			} else {
	//				key = match[1]
	//				//value := match[2]
	//				switch key {
	//				default:
	//					return nil, fmt.Errorf("unknown tag name: %s", key)
	//				}
	//			}
	//			if _, has := appliedKeys[key]; has {
	//				return nil, fmt.Errorf("duplicated tag: %s", key)
	//			}
	//			appliedKeys[key] = struct{}{}
	//		}
	//	}
	//}
	default:
		return nil, fmt.Errorf("unknown enumer generation kind: %s", generationKind)
	}
	structInfo.GenerateKind = generationKind

	return structInfo, nil
}
