package convert

import "strings"

// outOfScope 判断这句话是不是资产搜索表达不了的需求。
func outOfScope(text string) bool {
	if strings.Contains(text, "最安全") {
		return true
	}
	// 容器内存使用率依赖运行时监控，FOFA 没有对应字段。
	if strings.Contains(text, "内存使用率") {
		return true
	}
	return false
}
