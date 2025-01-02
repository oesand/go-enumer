package cases

type CaseType string

const (
	CamelCase  CaseType = "camelCase"
	KebabCase  CaseType = "kebab-case"
	PascalCase CaseType = "PascalCase"
	SnakeCase  CaseType = "snake_case"
)

func (ct CaseType) From(input string) string {
	if input == "" {
		return ""
	}
	switch ct {
	case CamelCase:
		return ToCamelCase(input)
	case KebabCase:
		return ToKebabCase(input)
	case PascalCase:
		return ToPascalCase(input)
	case SnakeCase:
		return ToSnakeCase(input)
	}
	return ""
}

func (ct CaseType) IsValid() bool {
	return ct == CamelCase ||
		ct == KebabCase ||
		ct == PascalCase ||
		ct == SnakeCase
}

func isAsciiUpperCase(r rune) bool {
	return 0x41 <= r && r <= 0x5a
}

func isAsciiLowerCase(r rune) bool {
	return 0x61 <= r && r <= 0x7a
}

func isAsciiDigit(r rune) bool {
	return 0x30 <= r && r <= 0x39
}

func toAsciiUpperCase(r rune) rune {
	if isAsciiLowerCase(r) {
		return r - 0x20
	}
	return r
}

func toAsciiLowerCase(r rune) rune {
	if isAsciiUpperCase(r) {
		return r + 0x20
	}
	return r
}
