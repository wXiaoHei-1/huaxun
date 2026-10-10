package convert

import (
	"strings"

	"huaxun/internal/fofa"
)

func matchForeign(text string) (fofa.Expr, bool) {
	switch {
	case strings.Contains(text, "Shodan"):
		return fofa.All(fofa.Eq("header", "Server: nginx"), fofa.Eq("title", "Welcome")), true
	case strings.Contains(text, "Google Dork"):
		return fofa.All(fofa.Eq("title", "Grafana"), fofa.Eq("body", "JavaScript")), true
	case strings.Contains(text, "urlscan.io"):
		return fofa.All(fofa.Eq("body", "login"), containsAny("host", "bet", "casino", "kasino")), true
	case strings.Contains(text, "Hunter"):
		return fofa.All(fofa.Eq("domain", "google"), fofa.Ne("domain", "google.com"), fofa.Exact("title", "")), true
	case strings.Contains(text, "Censys"):
		return fofa.Any(fofa.Exact("server", "Tengine"), fofa.Eq("header", "x-cache: MISS")), true
	case strings.Contains(text, "排除 ASN 4134"):
		// 香港是国家代码 HK。原句把 region 和全角引号写错了，按意图重写。
		return fofa.All(country("香港"), fofa.Exact("server", "Apache"), fofa.Ne("asn", "4134")), true
	case strings.Contains(text, "8080 端口仅状态码"):
		return fofa.All(
			fofa.Eq("app", "HFS"),
			fofa.Any(fofa.Eq("port", "80"), fofa.All(fofa.Eq("port", "8080"), fofa.Eq("status_code", "200"))),
		), true
	default:
		return nil, false
	}
}
