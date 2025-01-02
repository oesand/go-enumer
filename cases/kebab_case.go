package cases

// ToKebabCase converts a string to kebab-case.
func ToKebabCase(input string) string {
	if input == "" {
		return ""
	}

	runes := []rune(input)
	length := len(runes)
	result := make([]rune, 0, length)

	for i, r := range runes {
		if isAsciiUpperCase(r) {
			if i > 0 {
				result = append(result, '-')
			}
			result = append(result, toAsciiLowerCase(r))
		} else if isAsciiLowerCase(r) || isAsciiDigit(r) {
			result = append(result, r)
		} else {
			if len(result) > 0 && result[len(result)-1] != '-' {
				result = append(result, '-')
			}
		}
	}

	// Remove trailing hyphen if exists
	if len(result) > 0 && result[len(result)-1] == '-' {
		result = result[:len(result)-1]
	}

	return string(result)
}
