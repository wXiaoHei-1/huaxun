package fofa

import "strings"

type exprKind int

const (
	kindCmp exprKind = iota
	kindAnd
	kindOr
)

// Render 把表达式渲染成一条 FOFA 查询。
// 括号只加在并且、或者混用的位置，避免先且后或被写反。
func Render(expr Expr) string {
	text, _ := render(expr)
	return text
}

func render(expr Expr) (string, exprKind) {
	switch node := expr.(type) {
	case *Cmp:
		return renderCmp(node), kindCmp
	case *And:
		left, leftKind := render(node.Left)
		right, rightKind := render(node.Right)
		// 关键步骤：或者塞进并且时必须加括号。
		if leftKind == kindOr {
			left = "(" + left + ")"
		}
		if rightKind == kindOr {
			right = "(" + right + ")"
		}
		return left + " && " + right, kindAnd
	case *Or:
		left, leftKind := render(node.Left)
		right, rightKind := render(node.Right)
		// 关键步骤：并且塞进或者时必须加括号。
		if leftKind == kindAnd {
			left = "(" + left + ")"
		}
		if rightKind == kindAnd {
			right = "(" + right + ")"
		}
		return left + " || " + right, kindOr
	default:
		return "", kindCmp
	}
}

func renderCmp(cmp *Cmp) string {
	if cmp.Kind == ValueBool {
		return cmp.Field + "=" + cmp.Value
	}
	op := "="
	switch cmp.Op {
	case OpExact:
		op = "=="
	case OpNot:
		op = "!="
	case OpFuzzy:
		op = "*="
	}
	return cmp.Field + op + quote(cmp.Value)
}

// quote 把反斜杠和双引号转义后包进双引号。
func quote(value string) string {
	var b strings.Builder
	b.Grow(len(value) + 2)
	b.WriteByte('"')
	for i := 0; i < len(value); i++ {
		if value[i] == '\\' || value[i] == '"' {
			b.WriteByte('\\')
		}
		b.WriteByte(value[i])
	}
	b.WriteByte('"')
	return b.String()
}
