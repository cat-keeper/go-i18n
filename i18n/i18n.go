package i18n

import (
	"fmt"
	"reflect"
)

// I18n 表示一个语言实例
type I18n struct {
	Locale   string
	Messages map[string]string
	helper   *Helper
}

// T 获取国际化内容，并进行变量替换：支持 map 替换 ${key}，支持 {0},{1} 参数替换
func (i *I18n) T(key string, param interface{}) string {
	msg, ok := i.Messages[key]
	if !ok {
		return key
	}
	if param == nil {
		return msg
	}

	switch val := param.(type) {
	case map[string]string:
		tmp := make(map[string]interface{}, len(val))
		for k, v := range val {
			tmp[k] = v
		}
		return i.replaceWithMap(msg, tmp)
	case map[string]interface{}:
		return i.replaceWithMap(msg, val)
	default:
		// 尝试 slice/array 替换 {0},{1}
		rv := reflect.ValueOf(param)
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			var args []interface{}
			for i := 0; i < rv.Len(); i++ {
				args = append(args, rv.Index(i).Interface())
			}
			res, err := ReplaceArgs(msg, args...)
			if err == nil {
				return res
			}
		}
	}
	return msg
}

func (i *I18n) replaceWithMap(template string, data map[string]interface{}) string {
	resolver := func(placeholder string) (string, bool) {
		if v, ok := data[placeholder]; ok {
			return toString(v), true
		}
		return "", false
	}
	result, err := i.helper.Replace(template, resolver)
	if err != nil {
		return template
	}
	return result
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}
