package fofa

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

// 这些字段在一条资产上不能同时取两个精确值。
var exclusiveFields = map[string]bool{
	"country": true,
	"region":  true,
	"city":    true,
	"port":    true,
	"ip":      true,
	"asn":     true,
}

// Invalid 为真时，整题必须拒绝，不能渲染半截查询。
func Invalid(expr Expr) bool {
	if expr == nil {
		return true
	}
	bad := false
	walk(expr, func(cmp *Cmp) {
		// 关键步骤：端口越界或地址无法解析，整条表达式作废。
		if cmp.Field == "port" && cmp.Kind == ValueString && !ValidPort(cmp.Value) {
			bad = true
		}
		if cmp.Field == "ip" && cmp.Kind == ValueString && !ValidIPValue(cmp.Value) {
			bad = true
		}
	})
	return bad || Contradicts(expr)
}

// ValidPort 只接受 1 到 65535 的整数。
func ValidPort(value string) bool {
	port, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return port >= 1 && port <= 65535
}

// ValidIPValue 接受单个 IP 或 CIDR。
func ValidIPValue(value string) bool {
	if net.ParseIP(value) != nil {
		return true
	}
	_, _, err := net.ParseCIDR(value)
	return err == nil
}

// Contradicts 检查并且关系里必然同时成立的条件是否互相打架。
// 或者的两边不视为同时成立。标题里两个不同的词也不算矛盾。
func Contradicts(expr Expr) bool {
	grouped := map[string][]Cmp{}
	walkAnd(expr, func(cmp *Cmp) {
		if exclusiveFields[cmp.Field] {
			grouped[cmp.Field] = append(grouped[cmp.Field], *cmp)
		}
	})
	for _, items := range grouped {
		if fieldContradicts(items) {
			return true
		}
	}
	return false
}

func fieldContradicts(items []Cmp) bool {
	var exact []string
	notExact := map[string]bool{}
	for _, item := range items {
		switch item.Op {
		case OpExact:
			exact = append(exact, item.Value)
		case OpNot:
			notExact[item.Value] = true
		}
	}
	seen := map[string]bool{}
	for _, value := range exact {
		// 关键步骤：同一字段两个不同的精确值不能同时成立。
		if len(seen) > 0 && !seen[value] {
			return true
		}
		seen[value] = true
		if notExact[value] {
			return true
		}
	}
	return false
}

// IPv4BSegment 把一个 IPv4 收成它的 B 段。40.88.14.254 得到 40.88.0.0/16。
func IPv4BSegment(ip string) (string, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "", fmt.Errorf("无效 IP: %s", ip)
	}
	v4 := parsed.To4()
	if v4 == nil {
		return "", fmt.Errorf("不是 IPv4: %s", ip)
	}
	return fmt.Sprintf("%d.%d.0.0/16", v4[0], v4[1]), nil
}

// ExpandIPv4Closed 展开闭区间里的每一个 IPv4。区间超过 256 个地址时拒绝展开。
func ExpandIPv4Closed(start, end string) ([]string, error) {
	first, err := ipv4Uint(start)
	if err != nil {
		return nil, err
	}
	last, err := ipv4Uint(end)
	if err != nil {
		return nil, err
	}
	if last < first {
		return nil, fmt.Errorf("IP 区间起点大于终点")
	}
	if last-first > 256 {
		return nil, fmt.Errorf("IP 区间过大")
	}
	out := make([]string, 0, last-first+1)
	for value := first; ; value++ {
		out = append(out, formatIPv4(value))
		if value == last {
			break
		}
	}
	return out, nil
}

func ipv4Uint(ip string) (uint32, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return 0, fmt.Errorf("无效 IP: %s", ip)
	}
	v4 := parsed.To4()
	if v4 == nil {
		return 0, fmt.Errorf("不是 IPv4: %s", ip)
	}
	return binary.BigEndian.Uint32(v4), nil
}

func formatIPv4(value uint32) string {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, value)
	return net.IP(buf).String()
}

func walk(expr Expr, visit func(*Cmp)) {
	switch node := expr.(type) {
	case *Cmp:
		visit(node)
	case *And:
		walk(node.Left, visit)
		walk(node.Right, visit)
	case *Or:
		walk(node.Left, visit)
		walk(node.Right, visit)
	}
}

func walkAnd(expr Expr, visit func(*Cmp)) {
	switch node := expr.(type) {
	case *Cmp:
		visit(node)
	case *And:
		walkAnd(node.Left, visit)
		walkAnd(node.Right, visit)
	case *Or:
		// 或者的某一边可以不成立，不拿来判断整题矛盾。
	}
}
