package convert

import (
	"testing"

	"huaxun/internal/fofa"
)

func TestOpFromCue(t *testing.T) {
	cases := []struct {
		cue string
		op  fofa.Op
	}{
		{"是", fofa.OpExact},
		{"为", fofa.OpExact},
		{"等于", fofa.OpExact},
		{"完整为", fofa.OpExact},
		{"精确", fofa.OpExact},
		{"名为", fofa.OpExact},
		{"包含", fofa.OpContains},
		{"带有", fofa.OpContains},
		{"出现", fofa.OpContains},
		{"里有", fofa.OpContains},
		{"不是", fofa.OpNot},
		{"不包含", fofa.OpNot},
		{"排除", fofa.OpNot},
		{"不要", fofa.OpNot},
		{"恰好", fofa.OpFuzzy},
	}
	for _, tc := range cases {
		if opFromCue(tc.cue) != tc.op {
			t.Fatalf("说法 %s 的运算符不对", tc.cue)
		}
	}
}

func TestCueOnSentences(t *testing.T) {
	cases := []struct{ in, want string }{
		{"搜索网站名为 www.amazon.com 的资产。", `host=="www.amazon.com"`},
		{"我想找域名里带有 vpn 的所有网站。", `host="vpn"`},
		{"查 163.com 旗下的域名数据，要求它前面恰好还有两级子域，每一级都只有一个字符，例如 a.z.163.com。", `host*="?.?.163.com"`},
	}
	for _, tc := range cases {
		if got := ToQuery(tc.in); got != tc.want {
			t.Fatalf("in %s\ngot  %s\nwant %s", tc.in, got, tc.want)
		}
	}
}

func TestLogicHelpers(t *testing.T) {
	if got, want := fofa.Render(containsAny("title", "Sign in", "Login")), `title="Sign in" || title="Login"`; got != want {
		t.Fatalf("或: %s", got)
	}
	least := atLeastNonEmpty([]string{"a", "b", "c"}, 2)
	want := `(a!="" && b!="") || (a!="" && c!="") || (b!="" && c!="")`
	if got := fofa.Render(least); got != want {
		t.Fatalf("至少两项: %s", got)
	}
}

func TestLexicon(t *testing.T) {
	if countryCode("大不列颠及北爱尔兰联合王国") != "GB" || countryCode("英国") != "GB" {
		t.Fatal("英国代码应为 GB")
	}
	if countryCode("中国大陆") != "CN" || countryCode("香港") != "HK" {
		t.Fatal("中国大陆与香港代码不对")
	}
	if regions["江原道"] != "Gangwon-do" || cities["München"] != "Munich" {
		t.Fatal("地区或城市取值不对")
	}
	if clouds["亚马逊"] != "Amazon" || clouds["华为云"] != "HuaWeiCloud" {
		t.Fatal("云厂商取值不对")
	}
	if got := ToQuery("搜索韩国江原道的资产。"); got != `country="KR" && region="Gangwon-do"` {
		t.Fatalf("江原道: %s", got)
	}
}

func TestRejectAndForeign(t *testing.T) {
	reject := []string{
		"搜索开放 1000000 端口的资产。",
		"搜索地址或网段为 1.2.3.4.5 的资产。",
		"查同一条资产记录，国家代码必须等于 US，同时国家代码又必须不等于 US。",
		"帮我查最安全的网站。",
		"找出当前正在运行的容器内存使用率超过 80% 的服务器，没有容器监控数据或运行时访问权限。",
		"今天天气怎么样",
	}
	for _, text := range reject {
		if got := ToQuery(text); got != Rejected {
			t.Fatalf("应拒绝 %s，得到 %s", text, got)
		}
	}
	fixed := ToQuery("修正这条 FOFA 查询，意图是搜索香港地区、Server 精确为 Apache，并排除 ASN 4134：asn!=”4134″ && region=”HK” && server==”Apache”")
	if fixed != `country="HK" && server=="Apache" && asn!="4134"` {
		t.Fatalf("修正香港查询: %s", fixed)
	}
}
