package i18n

import (
	"strings"
)

// FindBestMatch 尝试从 locale 中回退获取最佳匹配语言
func FindBestMatch(locale string) string {
	if locale == "" {
		return defaultLocale
	}

	matches := generateFallbackChain(locale)

	for _, l := range matches {
		if _, ok := i18nInstances[l]; ok {
			return l
		}
	}

	if _, ok := i18nInstances[defaultLocale]; ok {
		return defaultLocale
	}
	return locale
}

// 生成 fallback 链：zh-Hans-CN → [zh-Hans-CN, zh-Hans, zh]
func generateFallbackChain(locale string) []string {
	segments := strings.Split(locale, "-")
	var chain []string
	for i := len(segments); i > 0; i-- {
		chain = append(chain, strings.Join(segments[:i], "-"))
	}
	return chain
}
