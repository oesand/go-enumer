package parse

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
)

var (
	tagsExp       = regexp.MustCompile(`\b(\w+)\s*:\s*(\w+)\b|\b\w+\b`)
	whitespaceExp = regexp.MustCompile(`\s+`)
)

func ensureTagHasValue(tag string, value string) error {
	if value == "" {
		return fmt.Errorf("tag %s should have a value", tag)
	}
	return nil
}

func visitAllTags(text string, preventDuplicate bool, cb func(string, string) error) (err error) {
	if text == "" {
		return nil
	}
	matches := tagsExp.FindAllStringSubmatch(text, -1)
	if matches != nil {
		var appliedKeys map[string]struct{}
		if preventDuplicate {
			appliedKeys = make(map[string]struct{})
		}
		for _, match := range matches {
			var key string
			if match[1] == "" && match[2] == "" {
				key = match[0]
				err = cb(key, "")
			} else {
				key = match[1]
				err = cb(key, match[2])
			}
			if preventDuplicate {
				if _, has := appliedKeys[key]; has {
					err = fmt.Errorf("duplicated tag: %s", key)
				} else {
					appliedKeys[key] = struct{}{}
				}
			}
			if err != nil {
				break
			}
		}
	}
	return
}

func newLocatedErr(fileSet *token.FileSet, fileName string, tspec *ast.TypeSpec, text string, a ...any) error {
	pos := tspec.Pos()
	file := fileSet.File(pos)
	line := file.Line(pos)
	return &locatedError{
		fileName: fileName,
		typeName: tspec.Name.Name,
		line:     line,
		text:     fmt.Sprintf(text, a...),
	}
}

type locatedError struct {
	fileName string
	typeName string
	line     int
	text     string
}

func (e *locatedError) Error() string {
	return fmt.Sprintf("[%s:%d @%s]: %s", e.fileName, e.line, e.typeName, e.text)
}
