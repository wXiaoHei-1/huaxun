package fofa

import "testing"

func TestQuoteAndBool(t *testing.T) {
	cases := []struct {
		name string
		expr Expr
		want string
	}{
		{name: "普通值", expr: Eq("title", "Welcome"), want: `title="Welcome"`},
		{name: "值内双引号", expr: Eq("body", `class="container"`), want: `body="class=\"container\""`},
		{name: "值内反斜杠", expr: Eq("body", `a\b`), want: `body="a\\b"`},
		{name: "布尔不加引号", expr: Bool("cert.is_valid", true), want: `cert.is_valid=true`},
		{name: "布尔假值", expr: Bool("cert.is_expired", false), want: `cert.is_expired=false`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Render(tc.expr); got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestParentheses(t *testing.T) {
	pureAnd := All(Eq("port", "80"), Eq("country", "CN"), Eq("app", "Grafana"))
	if got, want := Render(pureAnd), `port="80" && country="CN" && app="Grafana"`; got != want {
		t.Fatalf("纯且: got %s", got)
	}
	pureOr := Any(Eq("port", "80"), Eq("port", "443"))
	if got, want := Render(pureOr), `port="80" || port="443"`; got != want {
		t.Fatalf("纯或: got %s", got)
	}

	// HFS 的 80 全要，8080 仅在状态码 200 时要。
	hfs := All(
		Eq("app", "HFS"),
		Any(
			Eq("port", "80"),
			All(Eq("port", "8080"), Eq("status_code", "200")),
		),
	)
	want := `app="HFS" && (port="80" || (port="8080" && status_code="200"))`
	if got := Render(hfs); got != want {
		t.Fatalf("混合: got %s", got)
	}
}
