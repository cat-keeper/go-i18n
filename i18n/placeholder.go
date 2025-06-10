package i18n

import (
	"errors"
	"fmt"
	"strings"
)

// Helper 模仿SpringBoot的PropertyPlaceholderHelper写了go版本
type Helper struct {
	// 占位符前缀
	prefix string
	// 占位符后缀
	suffix string

	simplePrefix string
	// 默认值分隔符
	separator string
	// 是否忽略未解析的占位符
	ignoreUnresolvable bool
}

// Resolver 函数类型，用于解析占位符
type Resolver func(placeholder string) (string, bool)

// 已知的简单前缀映射
var wellKnownSimplePrefixes = map[string]string{
	"}": "{",
	"]": "[",
	")": "(",
}

// NewHelper 创建占位符解析助手
func NewHelper(prefix, suffix string) (*Helper, error) {
	return NewHelperWithOptions(prefix, suffix, "", true)
}

// NewHelperWithOptions 创建支持默认值分隔符的实例
func NewHelperWithOptions(prefix, suffix, separator string, ignoreUnresolvable bool) (*Helper, error) {
	if prefix == "" {
		return nil, errors.New("prefix must not be null")
	}
	if suffix == "" {
		return nil, errors.New("suffix must not be null")
	}

	simplePrefix := prefix
	if p, ok := wellKnownSimplePrefixes[suffix]; ok && strings.HasSuffix(prefix, p) {
		simplePrefix = p
	}

	return &Helper{
		prefix:             prefix,
		suffix:             suffix,
		simplePrefix:       simplePrefix,
		separator:          separator,
		ignoreUnresolvable: ignoreUnresolvable,
	}, nil
}

// Replace 替换字符串中的所有占位符
func (h *Helper) Replace(template string, resolver Resolver) (string, error) {
	if template == "" {
		return template, nil
	}
	return h.parse(template, resolver, nil)
}

// 递归解析字符串值
func (h *Helper) parse(template string, resolver Resolver, visitedPlaceholders map[string]bool) (string, error) {
	startIdx := strings.Index(template, h.prefix)
	if startIdx == -1 {
		return template, nil
	}

	var result strings.Builder
	result.WriteString(template[:startIdx])

	for startIdx != -1 {
		endIdx := h.findEndIndex(template, startIdx)
		if endIdx == -1 {
			break
		}

		placeholder := template[startIdx+len(h.prefix) : endIdx]
		originalPlaceholder := placeholder

		// 检测循环引用
		if visitedPlaceholders == nil {
			visitedPlaceholders = make(map[string]bool)
		}
		if visitedPlaceholders[originalPlaceholder] {
			return "", fmt.Errorf("circular placeholder reference '%s'", originalPlaceholder)
		}
		visitedPlaceholders[originalPlaceholder] = true

		// 递归解析占位符键中的占位符
		resolvedPlaceholder, err := h.parse(placeholder, resolver, visitedPlaceholders)
		if err != nil {
			return "", err
		}
		placeholder = resolvedPlaceholder

		// 解析占位符值
		propVal, ok := resolver(placeholder)
		if !ok && h.separator != "" {
			separatorIdx := strings.Index(placeholder, h.separator)
			if separatorIdx != -1 {
				actualPlaceholder := placeholder[:separatorIdx]
				defaultValue := placeholder[separatorIdx+len(h.separator):]
				propVal, ok = resolver(actualPlaceholder)
				if !ok {
					propVal = defaultValue
					ok = true
				}
			}
		}

		if ok {
			// 递归解析值中的占位符
			resolvedValue, err := h.parse(propVal, resolver, visitedPlaceholders)
			if err != nil {
				return "", err
			}
			result.WriteString(resolvedValue)
			nextStart := result.Len()
			result.WriteString(template[endIdx+len(h.suffix):])
			template = result.String()
			startIdx = strings.Index(template[nextStart:], h.prefix)
			if startIdx != -1 {
				startIdx += nextStart
			}
		} else if h.ignoreUnresolvable {
			result.WriteString(template[startIdx : endIdx+len(h.suffix)])
			startIdx = strings.Index(template[endIdx+len(h.suffix):], h.prefix)
			if startIdx != -1 {
				startIdx += endIdx + len(h.suffix)
			}
		} else {
			return "", fmt.Errorf("could not resolve placeholder '%s'", placeholder)
		}

		delete(visitedPlaceholders, originalPlaceholder)
	}

	return result.String(), nil
}

// 查找占位符结束位置，处理嵌套占位符
func (h *Helper) findEndIndex(s string, start int) int {
	index := start + len(h.prefix)
	nestedDepth := 0

	for index < len(s) {
		if strings.HasPrefix(s[index:], h.suffix) {
			if nestedDepth > 0 {
				nestedDepth--
				index += len(h.suffix)
			} else {
				return index
			}
		} else if strings.HasPrefix(s[index:], h.simplePrefix) {
			nestedDepth++
			index += len(h.simplePrefix)
		} else {
			index++
		}
	}
	return -1
}
