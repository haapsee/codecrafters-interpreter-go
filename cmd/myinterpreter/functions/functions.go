package functions

import (
	"fmt"
	"strings"
)

func TypeOf(obj interface{}) string {
	if obj == nil {
		return "nil"
	}

	switch obj.(type) {
	case bool:
		return "bool"
	case float64:
		return "float64"
	case string:
		return "string"
	}

	return ""
}

func IsAlphaDigit(char rune) bool {
	return IsAlpha(char) || IsDigit(char)
}

func IsAlpha(char rune) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || char == '_'
}

func IsDigit(char rune) bool {
	return strings.Contains("1234567890", string(char))
}

func FormatWithFixedPrecision(num float64) string {
	value := fmt.Sprintf("%g", num)

	if !strings.Contains(value, ".") {
		value += ".0"
	}

	return value
}
