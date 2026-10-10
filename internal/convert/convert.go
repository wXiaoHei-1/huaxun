// Package convert 把自然语言收成 FOFA 表达式。
// 每条规则只认自己的说法；没有任何规则命中时整题拒绝，不猜测查询。
package convert

import (
	"strings"

	"huaxun/internal/fofa"
)

// Rejected 是不能转换时的固定答案，不能改字，也不能附加解释。
const Rejected = "该需求不能直接转换为FOFA搜索语句"

// ToQuery 把一道题干变成一条查询或固定拒绝句。
func ToQuery(text string) string {
	text = strings.TrimSpace(text)
	if text == "" || outOfScope(text) {
		return Rejected
	}
	// 关键步骤：按固定顺序调用识别函数。锚点互不重叠，先命中的规则负责整句。
	matchers := []func(string) (fofa.Expr, bool){
		matchForeign,
		matchSpecial,
		matchCert,
		matchWeb,
		matchAsset,
		matchGeo,
		matchNet,
	}
	for _, match := range matchers {
		if expr, ok := match(text); ok {
			return finish(expr)
		}
	}
	return Rejected
}

func finish(expr fofa.Expr) string {
	// 关键步骤：端口、地址或自相矛盾不合法时，丢掉表达式，只回固定拒绝句。
	if expr == nil || fofa.Invalid(expr) {
		return Rejected
	}
	return fofa.Render(expr)
}
