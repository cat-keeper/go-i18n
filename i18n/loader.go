package i18n

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// LoadAllLocaleFiles 加载指定目录下所有符合 messages*.properties 的文件
func LoadAllLocaleFiles(dir string) map[string]map[string]string {
	locales := make(map[string]map[string]string)
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		name := info.Name()
		if strings.HasPrefix(name, "messages") && strings.HasSuffix(name, ".properties") {
			locale := extractLocale(name)
			props, err := loadProperties(path)
			if err == nil {
				locales[locale] = props
			}
		}
		return nil
	})
	return locales
}

// loadProperties 将 .properties 文件内容加载为 map
func loadProperties(filePath string) (map[string]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	props := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}
		idx := strings.IndexAny(line, "=:")

		if idx == -1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		props[key] = val
	}
	return props, nil
}

// extractLocale 从文件名提取 locale，如 messages_zh_TW.properties -> zh_TW
func extractLocale(filename string) string {
	base := strings.TrimSuffix(filename, ".properties")
	if base == "messages" {
		return "default"
	}
	if strings.HasPrefix(base, "messages_") {
		return strings.ReplaceAll(base[len("messages_"):], "-", "_")
	}
	return "default"
}
