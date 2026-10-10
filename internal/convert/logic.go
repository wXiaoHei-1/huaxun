package convert

import "huaxun/internal/fofa"

// opFromCue 把题干里的说法对应到 FOFA 运算符。
// 是、为、等于、完整为、精确、名为 → ==
// 包含、带有、出现 → =
// 不是、不包含、排除、不要 → !=
// 恰好（N 个字符或指定后缀形态）→ *=
func opFromCue(cue string) fofa.Op {
	switch cue {
	case "是", "为", "等于", "完整为", "精确", "名为":
		return fofa.OpExact
	case "不是", "不包含", "排除", "不要":
		return fofa.OpNot
	case "恰好":
		return fofa.OpFuzzy
	default:
		return fofa.OpContains
	}
}

func cmp(field, cue, value string) fofa.Expr {
	switch opFromCue(cue) {
	case fofa.OpExact:
		return fofa.Exact(field, value)
	case fofa.OpNot:
		return fofa.Ne(field, value)
	case fofa.OpFuzzy:
		return fofa.Fuzzy(field, value)
	default:
		return fofa.Eq(field, value)
	}
}

func containsAny(field string, values ...string) fofa.Expr {
	items := make([]fofa.Expr, 0, len(values))
	for _, value := range values {
		items = append(items, fofa.Eq(field, value))
	}
	return fofa.Any(items...)
}

func ports(values ...string) fofa.Expr {
	return containsAny("port", values...)
}

// atLeastNonEmpty 把「这几项里至少 n 项有值」展开成两两并且再取或。
func atLeastNonEmpty(fields []string, n int) fofa.Expr {
	var parts []fofa.Expr
	for _, group := range combinations(fields, n) {
		items := make([]fofa.Expr, 0, len(group))
		for _, field := range group {
			items = append(items, fofa.Ne(field, ""))
		}
		parts = append(parts, fofa.All(items...))
	}
	return fofa.Any(parts...)
}

func combinations(fields []string, n int) [][]string {
	var out [][]string
	var cur []string
	var walk func(start int)
	walk = func(start int) {
		if len(cur) == n {
			out = append(out, append([]string(nil), cur...))
			return
		}
		for i := start; i < len(fields); i++ {
			cur = append(cur, fields[i])
			walk(i + 1)
			cur = cur[:len(cur)-1]
		}
	}
	walk(0)
	return out
}
