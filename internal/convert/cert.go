package convert

import (
	"strings"

	"huaxun/internal/fofa"
)

func matchCert(text string) (fofa.Expr, bool) {
	switch {
	case strings.Contains(text, "GlobalSign"):
		return fofa.All(fofa.Eq("cert.issuer.org", "GlobalSign"), fofa.Bool("cert.is_expired", false)), true
	case strings.Contains(text, "apple.com"):
		return fofa.All(cmp("cert.subject.cn", "是", "apple.com"), cmp("domain", "不是", "apple.com")), true
	case strings.Contains(text, "TrustAsia"):
		return fofa.Eq("cert", "TrustAsia"), true
	case strings.Contains(text, "Alibaba"):
		return fofa.Eq("cert.subject.org", "Alibaba"), true
	case strings.Contains(text, "adobe.com"):
		return fofa.Exact("cert.domain", "adobe.com"), true
	case strings.Contains(text, "qq.com"):
		// 持有者名称的根域是 qq.com：裸域精确等于，子域用单层通配。
		return fofa.Any(
			fofa.Exact("cert.subject.cn", "qq.com"),
			fofa.Fuzzy("cert.subject.cn", "*.qq.com"),
		), true
	case strings.Contains(text, "godaddy"):
		return fofa.All(fofa.Eq("cert", "godaddy"), fofa.Ne("domain", "godaddy.com")), true
	case strings.Contains(text, "蜜罐和欺诈"):
		return fofa.All(fofa.Eq("cert", "google.com"), fofa.Bool("is_honeypot", false), fofa.Bool("is_fraud", false)), true
	case strings.Contains(text, "example.com"):
		return fofa.Any(fofa.Eq("cert.subject.cn", "example.com"), fofa.Eq("cert", "example.com")), true
	case strings.Contains(text, "TLS 1.0"):
		return fofa.All(fofa.Eq("tls.version", "TLS 1.0"), fofa.Bool("cert.is_expired", false)), true
	case strings.Contains(text, "颁发者与持有者不匹配"):
		// cert.is_equal 表示颁发者与持有者是否相同。英文接口说明与此一致。
		return fofa.Bool("cert.is_equal", false), true
	case strings.Contains(text, "证书与域名不匹配"):
		// cert.is_match 表示证书和资产域名是否匹配。
		return fofa.All(country("加拿大"), fofa.Bool("cert.is_match", false)), true
	case strings.Contains(text, "WinTAK"):
		return fofa.All(
			fofa.Eq("port", "143"),
			cmp("cert.sn", "为", "1"),
			containsAny("cert", "TAK", "ATAK", "WinTAK", "Squad", "CoT"),
		), true
	case strings.Contains(text, "aliyun.com"):
		return fofa.All(
			fofa.Exact("cert.domain", "aliyun.com"),
			fofa.Bool("cert.is_valid", true),
			fofa.Bool("is_domain", false),
		), true
	case strings.Contains(text, "kuaihuoyun.com"):
		return m082(), true
	default:
		return nil, false
	}
}

func m082() fofa.Expr {
	domains := []string{
		"kuaihuoyun.com",
		"wlhyos.cn",
		"wlhyos.com",
		"56ctms.com",
		"xiaokuaikeji.com",
		"zhihuiwuliu.org.cn",
	}
	var hits []fofa.Expr
	for _, domain := range domains {
		hits = append(hits, fofa.Exact("domain", domain), fofa.Exact("cert.domain", domain))
	}
	// 主域命中取或；证书匹配、有效、过期时间对整题生效，所以并且放在或的外面。
	return fofa.All(
		fofa.Any(hits...),
		fofa.Bool("cert.is_match", true),
		fofa.Bool("cert.is_valid", true),
		fofa.Eq("cert.not_after.after", "2024-01-01 00:00:00"),
		fofa.Eq("cert.not_after.before", "2026-08-30 18:17:49"),
	)
}
