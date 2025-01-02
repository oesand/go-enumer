package cases

// ToPascalCase converts a string to PascalCase.
func ToPascalCase(input string) string {
	if input == "" {
		return ""
	}

	runes := []rune(input)
	length := len(runes)
	result := make([]rune, 0, length)

	upperNext := true
	for _, r := range runes {
		if r == '_' || r == '-' || r == ' ' {
			upperNext = true
		} else {
			if upperNext {
				result = append(result, toAsciiUpperCase(r))
				upperNext = false
			} else {
				result = append(result, r)
			}
		}
	}

	return string(result)
}
