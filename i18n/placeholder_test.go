package i18n

import (
	"fmt"
	"testing"
)

func TestPlaceholder(t *testing.T) {
	// 基本示例
	basic, err := NewHelper("${", "}")
	if err != nil {
		t.Fatal(err)
	}
	actual := "Hello, Alice!"
	result, err := basic.Replace("Hello, ${name}!", func(s string) (string, bool) {
		if s == "name" {
			return "Alice", true
		}
		return "", false
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != actual {
		t.Errorf("got %s, want %s", result, actual)
	}
	fmt.Println(result) // Hello, Alice!

	// 带默认值
	withDefaults, err := NewHelperWithOptions("${", "}", ":", true)
	if err != nil {
		t.Fatal(err)
	}
	actual2 := "Color: blue"
	result2, err := withDefaults.Replace("Color: ${color:blue}", func(s string) (string, bool) {
		return "", false // 所有占位符都解析失败
	})
	if err != nil || result2 != actual2 {
		t.Errorf("got %s, want %s", result2, "Color: blue")
	}
	fmt.Println(result2) // Color: blue

	// 嵌套占位符
	actual3 := "nested"
	result3, err := basic.Replace("${outer}", func(s string) (string, bool) {
		if s == "outer" {
			return "${inner}", true
		}
		if s == "inner" {
			return "nested", true
		}
		return "", false
	})
	if err != nil {
		t.Fatal(err)
	}
	if result3 != actual3 {
		t.Errorf("got %s, want %s", result3, "nested")
	}
	fmt.Println(result3) // nested

	// 错误处理
	_, err = basic.Replace("${loop}", func(s string) (string, bool) {
		if s == "loop" {
			return "${loop}", true // 循环引用
		}
		return "", false
	})
	if err == nil {
		t.Errorf("got nil, want error")
	}
	fmt.Println(err) // 循环引用占位符: loop
}
