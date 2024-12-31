package parse

import (
	"fmt"
	"github.com/oesand/go-enumer/internal/shared"
	"regexp"
	"strings"
	"unicode"
)

var enumExp = regexp.MustCompile(`(?i)^\s*enum\(([^)]*)\)`)

func parseEnumType(typeName string, name string, comment string) (*shared.EnumInfo, error) {
	matches := enumExp.FindStringSubmatch(comment)
	if matches == nil {
		return nil, nil
	}
	valuesString := strings.ReplaceAll(matches[1], " ", "")
	if valuesString == "" {
		return nil, fmt.Errorf("empty enum values, see examples %s", shared.ProjectLink)
	}
	valueNames := strings.Split(valuesString, ",")

	var inverseNameOption bool
	prefixOption := name
	appliedKeys := make(map[string]struct{})

	enumEndIndex := enumExp.FindStringIndex(comment)[1]
	sequencedText := strings.Trim(comment[enumEndIndex:], " \n")
	if sequencedText != "" {
		matches := tagsExp.FindAllStringSubmatch(sequencedText, -1)
		if matches != nil {
			for _, match := range matches {
				var key string
				if match[1] == "" && match[2] == "" {
					key = match[0]
					switch key {
					case "inverse":
						inverseNameOption = true
					default:
						return nil, fmt.Errorf("unknown tag name: %s", match[0])
					}
				} else {
					key = match[1]
					value := match[2]
					switch key {
					case "prefix":
						prefixOption = toPascalCase(value)
					default:
						return nil, fmt.Errorf("unknown tag name: %s", key)
					}
				}
				if _, has := appliedKeys[key]; has {
					return nil, fmt.Errorf("duplicated tag: %s", key)
				}
				appliedKeys[key] = struct{}{}
			}
		}
	}

	values := make([]*shared.EnumValue, len(valueNames))

	for i, value := range valueNames {
		var name string
		if inverseNameOption {
			name = toPascalCase(value) + prefixOption
			if !unicode.IsLetter(rune(value[0])) {
				return nil, fmt.Errorf("generated invalid name for enum value(%s) with 'inverse' tag", value)
			}
		} else {
			name = prefixOption + toPascalCase(value)
		}

		values[i] = &shared.EnumValue{
			Name:  name,
			Value: value,
		}
	}

	enumInfo := &shared.EnumInfo{
		TypeName: shared.EnumSupportedTypes[typeName],
		EnumName: name,
		Values:   values,
	}

	return enumInfo, nil
}
