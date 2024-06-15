package utils

import (
	"strings"
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
