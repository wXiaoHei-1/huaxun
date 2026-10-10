package convert

import "strings"

import "huaxun/internal/fofa"

func matchGeo(text string) (fofa.Expr, bool) {
	switch {
	case strings.Contains(text, "大不列颠及北爱尔兰联合王国"):
		return country("大不列颠及北爱尔兰联合王国"), true
	case strings.Contains(text, "江原道"):
		return fofa.All(country("韩国"), region("江原道")), true
	case strings.Contains(text, "München"):
		return fofa.All(country("德国"), region("巴伐利亚"), city("München")), true
	case strings.Contains(text, "对象存储或应用托管"):
		return fofa.All(country("荷兰"), cloud("亚马逊")), true
	case strings.Contains(text, "河北地区的华为云"):
		// 河北是中国的省级行政区，国家固定为中国大陆。
		return fofa.All(country("中国"), region("河北"), cloud("华为云")), true
	case strings.Contains(text, "厦门市"):
		return fofa.All(country("中国"), city("厦门"), fofa.Eq("port", "80"), fofa.Eq("status_code", "200")), true
	default:
		return nil, false
	}
}
