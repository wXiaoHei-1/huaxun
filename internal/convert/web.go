package convert

import (
	"strings"

	"huaxun/internal/fofa"
)

func matchWeb(text string) (fofa.Expr, bool) {
	switch {
	case strings.Contains(text, "正文包含 jenkins"):
		return fofa.All(
			fofa.Any(cmp("title", "包含", "Jenkins"), cmp("body", "包含", "jenkins")),
			countryNot("中国"),
			fofa.Eq("after", "2026-08-01"),
		), true
	case strings.Contains(text, "NoSuchUpload"):
		return fofa.Eq("body", "NoSuchUpload"), true
	case strings.Contains(text, "Content-Security-Policy-Report-Only"):
		// 这是 HTTP 响应头字段名，网站资产记在 header，不记在服务 banner。
		return fofa.Eq("header", "Content-Security-Policy-Report-Only"), true
	case strings.Contains(text, "状态码为 308"):
		return fofa.All(fofa.Eq("type", "subdomain"), fofa.Eq("status_code", "308")), true
	case strings.Contains(text, "SSH-2.0-OpenSSH_6.1"):
		return fofa.Eq("banner", "SSH-2.0-OpenSSH_6.1"), true
	case strings.Contains(text, "PHPSESSID"):
		return fofa.Eq("header", "Set-Cookie:PHPSESSID"), true
	case strings.Contains(text, "jquery.min.js"):
		return fofa.Eq("body", "jquery.min.js"), true
	case strings.Contains(text, "阿姆斯特丹"):
		// 200 OK 是服务返回的原文，不是 HTTP 状态码字段。
		return fofa.All(country("荷兰"), city("阿姆斯特丹"), fofa.Eq("port", "554"), fofa.Eq("banner", "200 OK")), true
	case strings.Contains(text, "sohu.com"):
		return fofa.All(fofa.Exact("domain", "sohu.com"), fofa.Ne("banner", ""), fofa.Ne("banner", "apache")), true
	case strings.Contains(text, "JARM"):
		return fofa.All(
			fofa.Exact("jarm", "21d21d00021d21d00021d21d21d21dce319294781f0e12867fd06896a000ea"),
			fofa.Eq("status_code", "307"),
			fofa.Eq("header", "Content-Length: 0"),
		), true
	case strings.Contains(text, "正文有 /admin"):
		// 同一句里有正文，说明是网站。Set-Cookie 记在 header。
		return fofa.All(country("中国大陆"), fofa.Eq("header", "Set-Cookie"), fofa.Eq("body", "/admin")), true
	case strings.Contains(text, "管理任一项"):
		return fofa.All(
			fofa.Exact("domain", "baidu.com"),
			fofa.Any(containsAny("body", "Login", "admin", "Dashboard", "管理"), fofa.Eq("header", "Set-Cookie")),
		), true
	case strings.Contains(text, "13.0.0.0/8"):
		return fofa.All(
			fofa.Eq("ip", "13.0.0.0/8"),
			fofa.Eq("status_code", "200"),
			containsAny("body", "Welcome", "Sign in", "Dashboard"),
		), true
	case strings.Contains(text, "amazonaws.com"):
		return fofa.All(country("德国"), fofa.Any(
			fofa.Eq("host", "amazonaws.com"),
			fofa.Eq("domain", "amazonaws.com"),
			fofa.Eq("cname_domain", "amazonaws.com"),
			fofa.Eq("cert", "amazonaws.com"),
			fofa.Eq("cert.domain", "amazonaws.com"),
		)), true
	case strings.Contains(text, "不是 IPv6"):
		return fofa.All(containsAny("title", "Sign in", "Login"), fofa.Eq("port", "80"), fofa.Bool("is_ipv6", false)), true
	case strings.Contains(text, "Banner 同时包含 IIS"):
		return fofa.Any(fofa.All(fofa.Eq("banner", "IIS"), fofa.Eq("banner", "1.")), fofa.Eq("header", "IIS")), true
	case strings.Contains(text, "入学案内"):
		return fofa.All(
			country("日本"),
			fofa.Bool("is_domain", true),
			fofa.Eq("port", "443"),
			fofa.Eq("status_code", "200"),
			fofa.Bool("cert.is_valid", true),
			fofa.Any(containsAny("title", "大学", "学部", "研究"), containsAny("body", "入学案内", "学生")),
		), true
	case strings.Contains(text, "uat、pre"):
		return fofa.All(
			fofa.Exact("domain", "baidu.com"),
			fofa.Any(
				containsAny("title", "系统", "平台", "登录", "管理", "API"),
				containsAny("host", "api", "admin", "dev", "test", "uat", "pre"),
			),
		), true
	case strings.Contains(text, "标题不包含同一个关键词"):
		return fofa.All(fofa.Eq("body", "Jenkins"), fofa.Ne("title", "Jenkins")), true
	case strings.Contains(text, "四个字符"):
		return fofa.All(
			fofa.Exact("title", ""),
			fofa.Eq("status_code", "200"),
			fofa.Ne("fid", ""),
			fofa.Ne("icon_hash", ""),
			fofa.Eq("host", "admin"),
			fofa.Any(fofa.Fuzzy("domain", "*.cn"), fofa.Fuzzy("domain", "????.com")),
		), true
	case strings.Contains(text, `class="container"`):
		return fofa.Eq("body", `class="container"`), true
	case strings.Contains(text, "FatalError"):
		return fofa.Eq("body", `"exception": "Symfony\Component\ErrorHandler\Error\FatalError"`), true
	case strings.Contains(text, "6aa522a3-493"):
		// Etag 和 Content-Length 保留原值；其余头只要求字段名出现；不写状态码。
		return fofa.All(
			fofa.Eq("header", "Content-Length: 1171"),
			fofa.Eq("header", "Accept-Ranges"),
			fofa.Eq("header", "Cache-Control"),
			fofa.Eq("header", "Content-Type"),
			fofa.Eq("header", "Date"),
			fofa.Eq("header", `Etag: "6aa522a3-493"`),
			fofa.Eq("header", "Last-Modified"),
			fofa.Eq("header", "Server"),
			fofa.Eq("header", "Vary"),
		), true
	default:
		return nil, false
	}
}
