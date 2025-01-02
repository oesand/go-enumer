package cases

// ToSnakeCase converts a string to snake_case.
func ToSnakeCase(input string) string {
	if input == "" {
		return ""
	}

	runes := []rune(input)
	length := len(runes)
	result := make([]rune, 0, length)

	for i, r := range runes {
		if isAsciiUpperCase(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, toAsciiLowerCase(r))
		} else if isAsciiLowerCase(r) || isAsciiDigit(r) {
			result = append(result, r)
		} else {
			if len(result) > 0 && result[len(result)-1] != '_' {
				result = append(result, '_')
			}
		}
	}

	// Remove trailing underscore if exists
	if len(result) > 0 && result[len(result)-1] == '_' {
		result = result[:len(result)-1]
	}

	return string(result)
}
