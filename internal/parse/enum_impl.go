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

	enumEndIndex := enumExp.FindStringIndex(comment)[1]
	sequencedText := strings.Trim(comment[enumEndIndex:], " \n")
	if sequencedText != "" {
		err := visitAllTags(sequencedText, true, func(key, value string) (err error) {
			switch key {
			case "inverse":
				inverseNameOption = true
			case "prefix":
				err = ensureTagHasValue(key, value)
				prefixOption = cases.ToPascalCase(value)
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
		TypeName: shared.EnumSupportedTypes[typeName],
		EnumName: name,
		Values:   values,
	}

	return enumInfo, nil
}
