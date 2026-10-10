package convert

import (
	"strings"

	"huaxun/internal/fofa"
)

// jenkinsAgentFields 来自 Jenkins 文档 TCP Agent Listener 的诊断响应。
// 只保留字段名以产品名 Jenkins 或 Remoting 开头的项，不保留版本号和会话值。
// 示例里的 Client、Server 不以这两个前缀开头，不进入查询。
var jenkinsAgentFields = []string{
	"Jenkins-Agent-Protocols",
	"Jenkins-Version",
	"Jenkins-Session",
	"Remoting-Minimum-Version",
}

// octodexIconHash 是 https://octodex.github.com/favicon.ico 的 FOFA icon_hash。
// 算法：图标字节做标准 Base64，每 76 个字符换行并保留末尾换行，再取 MurmurHash3 32 位有符号值。
const octodexIconHash = "-1493334710"

func matchSpecial(text string) (fofa.Expr, bool) {
	switch {
	case strings.Contains(text, "TCP Agent Listener"):
		var hits []fofa.Expr
		for _, field := range jenkinsAgentFields {
			hits = append(hits, fofa.Eq("body", field), fofa.Eq("banner", field))
		}
		return fofa.Any(hits...), true
	case strings.Contains(text, "octodex.github.com"):
		return fofa.Eq("icon_hash", octodexIconHash), true
	case strings.Contains(text, "poliec.com"):
		// 未取得该仿冒页的标题、证书或图标，不能把域名本身当成仿冒特征。
		return nil, true
	case strings.Contains(text, "CVE-2021-41773"):
		// 只对应 Apache HTTP Server 2.4.49 这个产品版本，不写利用路径。
		return fofa.All(
			fofa.Exact("product", "Apache HTTP Server"),
			fofa.Exact("product.version", "2.4.49"),
		), true
	case strings.Contains(text, "ZeroSSL ECC DV SSL CA 2"):
		// 界面上的冒号分隔序列号收成 FOFA 可检索的连续大写十六进制。
		return fofa.All(
			fofa.Exact("cert.issuer.cn", "ZeroSSL ECC DV SSL CA 2"),
			fofa.Exact("cert.sn", "63BD83716214836EF80075356161EB06"),
		), true
	default:
		return nil, false
	}
}
