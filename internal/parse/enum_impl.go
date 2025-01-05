package parse

import (
	"fmt"
	"github.com/oesand/go-enumer/cases"
	"github.com/oesand/go-enumer/internal/shared"
	"regexp"
	"strings"
	"unicode"
)

var enumExp = regexp.MustCompile(`(?i)^\s*enum\(([^)]*)\)`)

func parseEnumType(enumType shared.KnownEnumType, name string, comment string) (*shared.EnumInfo, error) {
	matches := enumExp.FindStringSubmatch(comment)
	if matches == nil {
		return nil, nil
	}
	valuesString := whitespaceExp.ReplaceAllString(matches[1], "")
	if valuesString == "" {
		return nil, fmt.Errorf("empty enum values, see examples %s", shared.ProjectLink)
	}
	valueNames := strings.Split(valuesString, ",")

	var inverseNameOption bool
	prefixOption := name

	definedTags := make(map[string]string)
	enumEndIndex := enumExp.FindStringIndex(comment)[1]
	sequencedText := strings.Trim(comment[enumEndIndex:], " \n\t")
	if sequencedText != "" {
		err := visitAllTags(sequencedText, true, func(key, value string) (err error) {
			switch key {
			case "inverse":
				inverseNameOption = true
			case "prefix":
				err = ensureTagHasValue(key, value)
				prefixOption = cases.ToPascalCase(value)
			case "combined":
				if enumType != shared.IntEnum {
					return fmt.Errorf("tag \"combined\" is only allowed for int enum type")
				}
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

	values := make([]*shared.EnumValue, len(valueNames))
	for i, value := range valueNames {
		if value == "_" {
			if enumType != shared.IntEnum {
				return nil, fmt.Errorf("underscore value is only allowed for int enum type")
			}
			if i == 0 {
				return nil, fmt.Errorf("underscore value is not allowed as first value")
			}
			values[i] = nil
			continue
		}

		var name string
		if inverseNameOption {
			name = cases.ToPascalCase(value) + prefixOption
			if !unicode.IsLetter(rune(value[0])) {
				return nil, fmt.Errorf("generated invalid name for enum value(%s) with 'inverse' tag", value)
			}
		} else {
			name = prefixOption + cases.ToPascalCase(value)
		}

		values[i] = &shared.EnumValue{
			Name:  name,
			Value: value,
		}
	}

	enumInfo := &shared.EnumInfo{
		TypeName: enumType,
		EnumName: name,
		Values:   values,
		Tags:     definedTags,
	}

	return enumInfo, nil
}
