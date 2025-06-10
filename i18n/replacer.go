package i18n

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	// 支持 {0}, {1} 风格
	indexPattern = regexp.MustCompile(`\{(\d+)}`)
)

// ReplaceArgs 支持变长参数：替换 {0}, {1}
func ReplaceArgs(template string, args ...interface{}) (string, error) {
	if template == "" {
		return "", fmt.Errorf("empty template")
	}
	if len(args) == 0 {
		return template, nil
	}
	return indexPattern.ReplaceAllStringFunc(template, func(m string) string {
		indexStr := indexPattern.FindStringSubmatch(m)[1]
		index, err := strconv.Atoi(indexStr)
		if err != nil || index >= len(args) {
			return m
		}
		return fmt.Sprintf("%v", args[index])
	}), nil
}
