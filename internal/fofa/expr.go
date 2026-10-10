// Package fofa 只负责查询表达式、渲染和合法性。
// 它不理解中文题干。
package fofa

// Op 是 FOFA 允许的比较符。不在这里的符号不能渲染出去。
type Op int

const (
	// OpContains 对应 =，表示包含。
	OpContains Op = iota
	// OpExact 对应 ==，表示精确等于。
	OpExact
	// OpNot 对应 !=，表示不包含。否定只挂在单个比较上。
	OpNot
	// OpFuzzy 对应 *=，* 匹配任意长度，? 匹配一个字符。
	OpFuzzy
)

// ValueKind 区分带引号的文本和裸布尔值。
type ValueKind int

const (
	ValueString ValueKind = iota
	ValueBool
)

// Expr 是唯一的中间结构。只有比较、并且、或者三种节点。
type Expr interface {
	node()
}

// Cmp 是一条字段比较，例如 port="80"。
type Cmp struct {
	Field string
	Op    Op
	Value string
	Kind  ValueKind
}

func (c *Cmp) node() {}

// And 表示左右都要成立。
type And struct {
	Left  Expr
	Right Expr
}

func (a *And) node() {}

// Or 表示左右成立一个即可。
type Or struct {
	Left  Expr
	Right Expr
}

func (o *Or) node() {}

// Eq 构造包含比较。
func Eq(field, value string) Expr {
	return &Cmp{Field: field, Op: OpContains, Value: value, Kind: ValueString}
}

// Exact 构造精确比较。
func Exact(field, value string) Expr {
	return &Cmp{Field: field, Op: OpExact, Value: value, Kind: ValueString}
}

// Ne 构造否定比较。
func Ne(field, value string) Expr {
	return &Cmp{Field: field, Op: OpNot, Value: value, Kind: ValueString}
}

// Fuzzy 构造通配比较。
func Fuzzy(field, value string) Expr {
	return &Cmp{Field: field, Op: OpFuzzy, Value: value, Kind: ValueString}
}

// Bool 构造布尔比较。FOFA 的布尔值不加引号，固定用 =。
func Bool(field string, value bool) Expr {
	text := "false"
	if value {
		text = "true"
	}
	return &Cmp{Field: field, Op: OpContains, Value: text, Kind: ValueBool}
}

// All 把多条表达式用并且串起来。只有一条时不再包一层。
func All(items ...Expr) Expr {
	return fold(items, func(left, right Expr) Expr {
		return &And{Left: left, Right: right}
	})
}

// Any 把多条表达式用或者串起来。只有一条时不再包一层。
func Any(items ...Expr) Expr {
	return fold(items, func(left, right Expr) Expr {
		return &Or{Left: left, Right: right}
	})
}

func fold(items []Expr, join func(left, right Expr) Expr) Expr {
	var kept []Expr
	for _, item := range items {
		if item != nil {
			kept = append(kept, item)
		}
	}
	if len(kept) == 0 {
		return nil
	}
	expr := kept[0]
	for _, item := range kept[1:] {
		expr = join(expr, item)
	}
	return expr
}
