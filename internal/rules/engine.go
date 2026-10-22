// Package rules 实现 NL→FOFA 的规则引擎：
// 不可转换检测器 + 参数化通用规则 + 人工沉淀知识库，按优先级组成责任链，首中即停。
package rules

// FixedAnswer 比赛规定的固定答案文案，必须逐字精确（README 二、答案内容）。
const FixedAnswer = "该需求不能直接转换为FOFA搜索语句"

// Rule 一条转换规则。Convert 命中返回 (query, true)。
type Rule struct {
	Name     string
	Priority int                     // 数值越小越先执行
	Convert  func(text string) (string, bool)
}

// Apply 按优先级执行规则链，返回首个命中的结果与命中规则名。
func Apply(text string, chain []Rule) (query, rule string, ok bool) {
	for _, r := range chain {
		if q, hit := r.Convert(text); hit {
			return q, r.Name, true
		}
	}
	return "", "", false
}
