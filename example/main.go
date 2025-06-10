package main

import (
	"fmt"
	"github.com/cat-keeper/go-i18n/i18n"
)

func main() {
	// 设置默认目录（可省略）
	i18n.SetI18nDir("resources/i18n")

	// 普通使用：使用 map 替换 ${var}
	msg1 := i18n.Get("嗷嗷").T("test.message.01", map[string]interface{}{"name": "张三"})
	fmt.Println("zh:", msg1)

	// 使用位置参数替换 {0}, {1}
	msg2 := i18n.Get("zh_TW").T("test.message.02", []string{"end...", "hhhhhhh"})
	fmt.Println("zh_TW:", msg2)
}
