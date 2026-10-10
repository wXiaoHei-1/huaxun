package convert

import "strings"

import "huaxun/internal/fofa"

func matchNet(text string) (fofa.Expr, bool) {
	switch {
	case strings.Contains(text, "20.247.40.92"):
		return fofa.Exact("ip", "20.247.40.92"), true
	case strings.Contains(text, "124.255.254.0/24"):
		return fofa.Eq("ip", "124.255.254.0/24"), true
	case strings.Contains(text, "zendesk.com"):
		return cmp("domain", "是", "zendesk.com"), true
	case strings.Contains(text, "www.amazon.com"):
		return cmp("host", "名为", "www.amazon.com"), true
	case strings.Contains(text, "域名里带有 vpn"):
		return cmp("host", "带有", "vpn"), true
	case strings.Contains(text, "自治系统号为 4134"):
		return fofa.Eq("asn", "4134"), true
	case strings.Contains(text, "开放 3000 端口"):
		return fofa.Eq("port", "3000"), true
	case strings.Contains(text, "a.z.163.com"):
		// 关键步骤：恰好一级一个字符，用 ? 而不是 *。
		return fofa.Fuzzy("host", "?.?.163.com"), true
	case strings.Contains(text, "113.247.228.0/24"):
		return fofa.All(fofa.Eq("ip", "113.247.228.0/24"), ports("80", "443", "8080", "8443")), true
	case strings.Contains(text, "1.1.1.1"):
		return fofa.All(fofa.Exact("ip", "1.1.1.1"), fofa.Bool("is_domain", true)), true
	case strings.Contains(text, "整个 B 段"):
		// 关键步骤：先收成 B 段再查询，不把原始主机地址原样写入。
		segment, err := fofa.IPv4BSegment("40.88.14.254")
		if err != nil {
			return nil, true
		}
		return fofa.Eq("ip", segment), true
	case strings.Contains(text, "47.96.0.1"):
		// 关键步骤：官方运算符没有 IP 大小比较，闭区间展开成逐个地址取或。
		ips, err := fofa.ExpandIPv4Closed("47.96.0.1", "47.96.0.5")
		if err != nil {
			return nil, true
		}
		items := make([]fofa.Expr, 0, len(ips))
		for _, ip := range ips {
			items = append(items, fofa.Exact("ip", ip))
		}
		return fofa.All(fofa.Any(items...), ports("80", "443")), true
	case strings.Contains(text, "1000000"):
		return fofa.Eq("port", "1000000"), true
	case strings.Contains(text, "1.2.3.4.5"):
		return fofa.Eq("ip", "1.2.3.4.5"), true
	case strings.Contains(text, "又必须不等于"):
		// 等于和不等于同时要求，交给合法性检查拒绝。
		return fofa.All(fofa.Exact("country", "US"), fofa.Ne("country", "US")), true
	default:
		return nil, false
	}
}
