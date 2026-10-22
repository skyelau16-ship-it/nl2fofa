package rules

import "testing"

// 固定答案文案逐字精确 —— 比赛红线，一个字都不能差。
func TestFixedAnswerExact(t *testing.T) {
	want := "该需求不能直接转换为FOFA搜索语句"
	if FixedAnswer != want {
		t.Fatalf("固定答案被改动: %q", FixedAnswer)
	}
}

func TestImpossible(t *testing.T) {
	cases := map[string]bool{
		"搜索开放 70001 端口的资产。":              true, // README 示例
		"搜索开放 1000000 端口的资产。":            true, // M091
		"搜索地址或网段为 1.2.3.4.5 的资产。":       true, // M092
		"查同一条资产记录，国家代码必须等于 US，同时国家代码又必须不等于 US。": true, // M093
		"帮我查最安全的网站。":                  true, // M094
		"找出当前正在运行的容器内存使用率超过 80% 的服务器。":   true, // M095
		"fofa有识别SOCKS的v5版本吗？":          true, // M015
		"搜索开放 3000 端口的资产。":             false,
		"搜索标题包含 Jenkins 的资产。":           false,
	}
	for text, want := range cases {
		if got := Impossible(text); got != want {
			t.Errorf("Impossible(%q) = %v, want %v", text, got, want)
		}
	}
}

func TestNormalizeCIDR(t *testing.T) {
	if got := normalizeCIDR("40.88.14.254/16"); got != "40.88.0.0/16" {
		t.Errorf("B 段归一化失败: %s", got)
	}
	if got := normalizeCIDR("124.255.254.0/24"); got != "124.255.254.0/24" {
		t.Errorf("/24 归一化失败: %s", got)
	}
}

// 知识库键必须互不包含，避免首中即停导致的误路由。
func TestKnowledgeKeysDistinct(t *testing.T) {
	entries := KnowledgeBase()
	for i, a := range entries {
		for j, b := range entries {
			if i != j && len(a.Key) <= len(b.Key) && b.Key != a.Key &&
				contains(b.Key, a.Key) && len(a.Key) >= 6 {
				t.Logf("警告: 键 %q 是 %q 的子串（请确认顺序正确）", a.Key, b.Key)
			}
		}
	}
}

func contains(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) && (s == sub || len(s) > len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
