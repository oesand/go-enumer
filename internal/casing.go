package internal

func toPascalCase(input string) string {
	result := make([]rune, 0, len(input))

	const (
		ChIsFirstOfStr = iota
		ChIsNextOfUpper
		ChIsNextOfMark
		ChIsOthers
	)
	var flag uint8 = ChIsFirstOfStr

	for _, ch := range input {
		if isAsciiUpperCase(ch) {
			switch flag {
			case ChIsNextOfUpper:
				result = append(result, toAsciiLowerCase(ch))
				//flag = ChIsNextOfUpper
			default:
				result = append(result, ch)
				flag = ChIsNextOfUpper
			}
		} else if isAsciiLowerCase(ch) {
			switch flag {
			case ChIsNextOfUpper:
				n := len(result)
				prev := result[n-1]
				if isAsciiLowerCase(prev) {
					result[n-1] = toAsciiUpperCase(prev)
				}
				result = append(result, ch)
				flag = ChIsOthers
			case ChIsFirstOfStr, ChIsNextOfMark:
				result = append(result, toAsciiUpperCase(ch))
				flag = ChIsNextOfUpper
			default:
				result = append(result, ch)
				flag = ChIsOthers
			}
		} else if isAsciiDigit(ch) {
			result = append(result, ch)
			flag = ChIsNextOfMark
		} else {
			if flag != ChIsFirstOfStr {
				flag = ChIsNextOfMark
			}
		}
	}

	return string(result)
}

func isAsciiUpperCase(r rune) bool {
	return (0x41 <= r && r <= 0x5a)
}

func isAsciiLowerCase(r rune) bool {
	return (0x61 <= r && r <= 0x7a)
}

func isAsciiDigit(r rune) bool {
	return (0x30 <= r && r <= 0x39)
}

func toAsciiUpperCase(r rune) rune {
	return (r + 0x41 - 0x61)
}

func toAsciiLowerCase(r rune) rune {
	return (r + 0x61 - 0x41)
}
