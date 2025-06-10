package i18n

import (
	"sync"
)

var (
	i18nInstances = make(map[string]*I18n)
	defaultLocale = "default"
	i18nDir       = "./resource/i18n/"
	initOnce      sync.Once
	mu            sync.RWMutex
)

// InitMessagesFile 初始化默认目录下的所有语言文件
func InitMessagesFile() {
	initOnce.Do(func() {
		files := LoadAllLocaleFiles(i18nDir)
		for locale, messages := range files {
			h, _ := NewHelper("${", "}")
			i18nInstances[locale] = &I18n{
				Locale:   locale,
				Messages: messages,
				helper:   h,
			}
		}
	})
}

// Get 返回指定语言的实例，若未找到则根据 fallback 链尝试
func Get(locale string) *I18n {
	mu.RLock()
	defer mu.RUnlock()

	bestMatch := FindBestMatch(locale)
	if inst, ok := i18nInstances[bestMatch]; ok {
		return inst
	}
	return &I18n{Locale: locale, Messages: map[string]string{}}
}

// SetI18nDir 修改资源目录路径（修改后需重新调用 InitMessagesFile）
func SetI18nDir(dir string) {
	mu.Lock()
	defer mu.Unlock()
	i18nDir = dir
	initOnce = sync.Once{} // 重置 once，允许重新初始化
	InitMessagesFile()
}

func init() {
	InitMessagesFile()
}
