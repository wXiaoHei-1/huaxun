package fofa

import "testing"

func TestValidPort(t *testing.T) {
	ok := []string{"1", "80", "65535"}
	bad := []string{"0", "65536", "1000000", "80a", "", "-1"}
	for _, value := range ok {
		if !ValidPort(value) {
			t.Fatalf("应接受端口 %s", value)
		}
	}
	for _, value := range bad {
		if ValidPort(value) {
			t.Fatalf("应拒绝端口 %s", value)
		}
	}
	if !Invalid(Eq("port", "1000000")) {
		t.Fatal("越界端口应使整条表达式无效")
	}
}

func TestValidIPValue(t *testing.T) {
	if !ValidIPValue("20.247.40.92") || !ValidIPValue("124.255.254.0/24") {
		t.Fatal("合法地址被拒绝")
	}
	if ValidIPValue("1.2.3.4.5") || !Invalid(Eq("ip", "1.2.3.4.5")) {
		t.Fatal("非法地址应被拒绝")
	}
}

func TestBSegmentAndRange(t *testing.T) {
	got, err := IPv4BSegment("40.88.14.254")
	if err != nil || got != "40.88.0.0/16" {
		t.Fatalf("B 段 got %s err %v", got, err)
	}
	if _, err := IPv4BSegment("1.2.3.4.5"); err == nil {
		t.Fatal("非法 IP 不应算出 B 段")
	}
	ips, err := ExpandIPv4Closed("47.96.0.1", "47.96.0.5")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"47.96.0.1", "47.96.0.2", "47.96.0.3", "47.96.0.4", "47.96.0.5"}
	if len(ips) != len(want) {
		t.Fatalf("got %v", ips)
	}
	for i := range want {
		if ips[i] != want[i] {
			t.Fatalf("got %v", ips)
		}
	}
	if _, err := ExpandIPv4Closed("47.96.0.5", "47.96.0.1"); err == nil {
		t.Fatal("反向区间应失败")
	}
}

func TestContradiction(t *testing.T) {
	clash := All(Exact("country", "US"), Ne("country", "US"))
	if !Contradicts(clash) || !Invalid(clash) {
		t.Fatal("国家代码既等于又不等于 US 应拒绝")
	}
	twoCountries := All(Exact("country", "US"), Exact("country", "CN"))
	if !Contradicts(twoCountries) {
		t.Fatal("两个精确国家应拒绝")
	}
	titles := All(Exact("title", "Jenkins"), Exact("title", "Mail"))
	if Contradicts(titles) || Invalid(titles) {
		t.Fatal("标题里的两个词不是矛盾")
	}
	either := Any(Exact("country", "US"), Ne("country", "US"))
	if Contradicts(either) {
		t.Fatal("或者的两边不必同时成立")
	}
}
