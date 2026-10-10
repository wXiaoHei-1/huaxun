package convert

import "huaxun/internal/fofa"

// 国家代码用 ISO。中国大陆在 FOFA 里就是 CN，港澳台分别是 HK、MO、TW。
// 因此「中国大陆」「不含港澳台的中国」只写 country="CN"，不再叠加地区排除。
var countries = map[string]string{
	"中国":   "CN",
	"中国大陆": "CN",
	"香港":   "HK",
	"澳门":   "MO",
	"台湾":   "TW",
	"美国":   "US",
	"日本":   "JP",
	"韩国":   "KR",
	"德国":   "DE",
	"荷兰":   "NL",
	"英国":   "GB",
	"大不列颠及北爱尔兰联合王国": "GB",
	"加拿大": "CA",
	"土耳其": "TR",
	"新西兰": "NZ",
}

// 地区、城市对照 FOFA 结果页的写法（如 Germany / Hessen，Netherlands / Noord-Holland / Amsterdam）。
// 德国州名用德文 Bayern，与 Hessen 一致；München 用 IP 库里的英文城市名 Munich。
// 江原道用 Gangwon-do。河北用 Hebei。
var regions = map[string]string{
	"江原道":  "Gangwon-do",
	"巴伐利亚": "Bayern",
	"河北":   "Hebei",
}

var cities = map[string]string{
	"München": "Munich",
	"厦门":      "Xiamen",
	"阿姆斯特丹":   "Amsterdam",
}

// 云厂商取值来自公开整理的 cloud_name：亚马逊是 Amazon，华为云是 HuaWeiCloud。
// DigitalOcean 不在那份清单里，按同一字段的厂商名写入。
var clouds = map[string]string{
	"亚马逊":          "Amazon",
	"华为云":          "HuaWeiCloud",
	"DigitalOcean": "DigitalOcean",
}

func countryCode(name string) string {
	code, ok := countries[name]
	if !ok {
		panic("未登记的国家: " + name)
	}
	return code
}

func country(name string) fofa.Expr {
	return fofa.Eq("country", countryCode(name))
}

func countryNot(name string) fofa.Expr {
	return fofa.Ne("country", countryCode(name))
}

func region(name string) fofa.Expr {
	value, ok := regions[name]
	if !ok {
		panic("未登记的地区: " + name)
	}
	return fofa.Eq("region", value)
}

func city(name string) fofa.Expr {
	value, ok := cities[name]
	if !ok {
		panic("未登记的城市: " + name)
	}
	return fofa.Eq("city", value)
}

func cloud(name string) fofa.Expr {
	value, ok := clouds[name]
	if !ok {
		panic("未登记的云厂商: " + name)
	}
	return fofa.Eq("cloud_name", value)
}
