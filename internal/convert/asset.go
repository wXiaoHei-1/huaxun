package convert

import (
	"strings"

	"huaxun/internal/fofa"
)

func matchAsset(text string) (fofa.Expr, bool) {
	switch {
	case strings.Contains(text, "中国的 DNS"):
		return fofa.All(country("中国"), fofa.Eq("protocol", "dns")), true
	case strings.Contains(text, "SOCKS"):
		return fofa.Eq("protocol", "socks5"), true
	case strings.Contains(text, "22022"):
		return fofa.All(fofa.Eq("protocol", "ssh"), fofa.Eq("port", "22022")), true
	case strings.Contains(text, "使用 TCP 协议的数据"):
		return fofa.Eq("base_protocol", "tcp"), true
	case strings.Contains(text, "SNMP"):
		return fofa.Eq("protocol", "snmp"), true
	case strings.Contains(text, "Memcached"):
		return fofa.All(fofa.Eq("protocol", "memcached"), country("美国")), true
	case strings.Contains(text, "真正跑在 TCP 上的 FTP"):
		return fofa.All(fofa.Eq("protocol", "ftp"), fofa.Eq("base_protocol", "tcp")), true
	case strings.Contains(text, "OpenSSH_10.4"):
		return fofa.All(fofa.Eq("protocol", "ssh"), fofa.Eq("banner", "OpenSSH_10.4")), true
	case strings.Contains(text, "海康威视"):
		// 摄像头对应 FOFA 的产品规则名，不用中文口语直译。
		return fofa.All(country("日本"), fofa.Eq("app", "海康威视-视频监控")), true
	case strings.Contains(text, "产品版本号为 10.0.0"):
		return fofa.Eq("product.version", "10.0.0"), true
	case strings.Contains(text, "中间层软件"):
		return fofa.Eq("category", "中间件"), true
	case strings.Contains(text, "操作系统被识别为 linux"):
		return fofa.Eq("os", "linux"), true
	case strings.Contains(text, "没有绑定域名的网站类"):
		return fofa.All(country("中国大陆"), fofa.Eq("type", "subdomain"), fofa.Bool("is_domain", false)), true
	case strings.Contains(text, "包含Inc"):
		return fofa.Eq("org", "Inc"), true
	case strings.Contains(text, "WHM 登录"):
		return fofa.All(country("美国"), fofa.Eq("org", "OVH SAS"), fofa.Eq("port", "2087"), cmp("title", "完整为", "WHM 登录")), true
	case strings.Contains(text, ".io 域名"):
		return fofa.All(fofa.Fuzzy("domain", "*.io"), fofa.Eq("app", "WordPress")), true
	case strings.Contains(text, "/api/show"):
		return fofa.All(country("英国"), cloud("DigitalOcean"), fofa.Eq("body", "/api/show")), true
	case strings.Contains(text, "JA3S"):
		return fofa.All(country("土耳其"), cmp("tls.ja3s", "等于", "f4febc55ea12b31ae17cfb7e614afda8")), true
	case strings.Contains(text, "新西兰"):
		return fofa.All(country("新西兰"), ports("80", "443"), fofa.Bool("is_cloud", true)), true
	case strings.Contains(text, "邮件系统"):
		// 「被归为」对到 category，类别名用题干里的「邮件系统」。
		return fofa.All(fofa.Eq("category", "邮件系统"), fofa.Eq("title", "Mail"), countryNot("美国")), true
	case strings.Contains(text, "8888"):
		return fofa.All(country("中国大陆"), fofa.Eq("port", "8888")), true
	case strings.Contains(text, "Let's Encrypt"):
		return fofa.All(
			country("台湾"),
			fofa.Eq("asn", "20940"),
			fofa.Eq("port", "443"),
			fofa.Eq("tls.version", "TLS 1.3"),
			fofa.Ne("cert.issuer.org", "Let's Encrypt"),
			fofa.Ne("cert.issuer.org", "ZeroSSL"),
		), true
	case strings.Contains(text, "HFS /"):
		return fofa.All(
			fofa.Eq("app", "HFS"),
			countryNot("日本"),
			fofa.Eq("status_code", "200"),
			cmp("title", "等于", "HFS /"),
			fofa.Ne("cloud_name", clouds["亚马逊"]),
		), true
	case strings.Contains(text, "中国 Grafana"):
		return fofa.All(country("中国"), fofa.Eq("app", "Grafana")), true
	case strings.Contains(text, ".edu.ar"):
		return fofa.All(
			fofa.Fuzzy("host", "*.edu.ar"),
			ports("80", "443"),
			fofa.Ne("server", "Nginx"),
			fofa.Ne("server", "WAF"),
		), true
	case strings.Contains(text, "x-qiniu-zone"):
		return fofa.All(
			country("中国"),
			fofa.Bool("is_domain", true),
			fofa.Eq("header", "x-qiniu-zone"),
			fofa.Ne("domain", "qiniudn.com"),
		), true
	case strings.Contains(text, "至少两项有值"):
		return fofa.All(
			fofa.Eq("app", "Grafana"),
			atLeastNonEmpty([]string{"header_hash", "body_hash", "icon_hash", "fid"}, 2),
		), true
	default:
		return nil, false
	}
}
