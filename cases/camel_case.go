package cases

// ToCamelCase converts a string to camelCase.
func ToCamelCase(input string) string {
	if input == "" {
		return ""
	}

	runes := []rune(input)
	length := len(runes)
	result := make([]rune, 0, length)

	upperNext := false
	for i, r := range runes {
		if r == '_' || r == '-' || r == ' ' {
			upperNext = true
		} else {
			if upperNext {
				result = append(result, toAsciiUpperCase(r))
				upperNext = false
			} else {
				if i == 0 {
					result = append(result, toAsciiLowerCase(r))
				} else {
					result = append(result, r)
				}
			}
		}
	}

	return string(result)
}
