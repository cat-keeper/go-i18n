# Go-I18n Toolkit

一个受 SpringBoot MessageSource 启发的国际化工具包，使用 `.properties` 文件作为语言资源，支持语言回退、占位符替换等功能。

## ✨ 功能特点

- ✅ 兼容 SpringBoot `.properties` 文件格式
- ✅ 支持 `${key}` 与 `{0}` 样式变量替换
- ✅ 多语言 fallback：如 `en-US` → `en` → `default`
- ✅ 热插拔国际化目录配置

## 📦 使用示例

```go
// 设置默认目录（可省略）
i18n.SetI18nDir("resources/i18n")

// 普通使用：使用 map 替换 ${var}
msg1 := i18n.Get("嗷嗷").T("test.message.01", map[string]interface{}{"name": "张三"})
fmt.Println("zh:", msg1)

// 使用位置参数替换 {0}, {1}
msg2 := i18n.Get("zh_TW").T("test.message.02", []string{"end...", "hhhhhhh"})
fmt.Println("zh_TW:", msg2)
```

## 📦 默认国际化文件结构
``` properties
resource/
└── i18n/
    ├── messages.properties         # 默认
    ├── messages_zh.properties      # 简体中文
    ├── messages_en.properties      # 英文
    └── messages_zh_TW.properties   # 繁体中文
```

## 🛠 变更语言资源目录

```go
i18n.SetDir("./config/i18n/")
```

## 📜 LICENSE

MIT License - see [LICENSE](https://github.com/go-i18n/i18n/blob/master/LICENSE)

## 🤝 欢迎指正

go新手, 有什么不正确的地方欢迎指正 ❤️
