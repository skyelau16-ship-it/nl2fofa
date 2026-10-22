package validate

import "testing"

func TestValid(t *testing.T) {
	cases := []struct {
		q    string
		want bool
	}{
		{`app="gitlab"`, true},
		{`(title="Sign in" || title="Login") && port="80"`, true},
		{`host~"^[^.]\.[^.]\.163\.com$"`, true},
		{`body="class=\"container\""`, true},
		{`(title="A" && port="80"`, false},  // 括号失衡
		{`title="A`, false},                  // 引号不成对
		{`unknown_field="x"`, false},         // 未知字段
		{`ip.after="1.1.1.1" && ip.before="1.1.1.5"`, true},
	}
	for _, c := range cases {
		if got, _ := Valid(c.q); got != c.want {
			t.Errorf("Valid(%q) = %v, want %v", c.q, got, c.want)
		}
	}
}
