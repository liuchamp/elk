package utils

import (
	"strings"
	"unicode"
)

func ToCamelCase(s string) string {
	// 将字符串按下划线分隔成单词
	words := strings.Split(s, "_")

	// 遍历单词列表并转换成驼峰命名
	var result string
	for _, word := range words {
		// 将单词的第一个字母转换成大写
		result += strings.Title(word)
	}

	return result
}

// SnakeToCamel converts a snake_case string to camelCase
func SnakeToCamel(s string) string {
	// Split the string by underscores
	parts := strings.Split(s, "_")
	for i := 0; i < len(parts); i++ {
		// Convert each part to title case except the first part
		if i > 0 {
			parts[i] = strings.Title(parts[i])
		}
	}

	// Join the parts
	result := strings.Join(parts, "")

	// Ensure the first character is lowercase
	if len(result) > 0 {
		result = string(unicode.ToLower(rune(result[0]))) + result[1:]
	}

	return result
}

// PascalToCamel converts a PascalCase string to camelCase
func PascalToCamel(s string) string {
	if len(s) == 0 {
		return s
	}

	// Convert the first character to lowercase
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])

	return string(runes)
}
