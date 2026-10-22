// Package pipeline 编排转换流水线：归一化 → 不可转换检测 → 规则链 → LLM 兜底 → 语法校验。
// 最终以答案模板为骨架填充，保证参赛包编号与题号天然一致（README 红线）。
package pipeline

import (
	"fmt"
	"strings"

	"github.com/skyelau16-ship-it/nl2fofa/internal/llm"
	"github.com/skyelau16-ship-it/nl2fofa/internal/normalize"
	"github.com/skyelau16-ship-it/nl2fofa/internal/rules"
	"github.com/skyelau16-ship-it/nl2fofa/internal/validate"
)

// Question 题目.json 的单条记录。
type Question struct {
	ID         string `json:"题号"`
	NaturalLan string `json:"自然语言输入"`
}

// AnswerItem 答案数组的单条记录。
type AnswerItem struct {
	ID    string `json:"题号"`
	Query string `json:"查询语句"`
}

// Template 答案模板.json 的结构，输出沿用同一结构。
type Template struct {
	Player  string       `json:"选手名称"`
	PkgID   string       `json:"参赛包编号"`
	Answers []AnswerItem `json:"答案"`
}

// Stats 各来源命中统计。
type Stats struct {
	Impossible, Knowledge, General, LLM, LLMFallbackFixed int
}

// Run 对全部题目执行流水线，返回按模板顺序填充的答案与统计。
func Run(questions []Question, tpl Template, player string, client *llm.Client) (Template, Stats) {
	byID := make(map[string]string, len(questions))
	for _, q := range questions {
		byID[q.ID] = q.NaturalLan
	}
	var st Stats
	chain := rules.Chain()
	for i := range tpl.Answers {
		id := tpl.Answers[i].ID
		text, exists := byID[id]
		if !exists {
			tpl.Answers[i].Query = rules.FixedAnswer // 防御：模板题号缺题目时兜底
			continue
		}
		query, source := convert(text, chain, client, &st)
		_ = source
		tpl.Answers[i].Query = query
	}
	tpl.Player = player
	return tpl, st
}

// convert 单题转换：规则优先，LLM 兜底，校验失败降级为固定答案。
func convert(text string, chain []rules.Rule, client *llm.Client, st *Stats) (string, string) {
	norm := normalize.Text(text)

	if rules.Impossible(norm) {
		st.Impossible++
		return rules.FixedAnswer, "impossible"
	}
	if q, ruleName, ok := rules.Apply(norm, chain); ok {
		if strings.HasPrefix(ruleName, "knowledge") {
			st.Knowledge++
		} else {
			st.General++
		}
		if ok2, _ := validate.Valid(q); ok2 {
			return q, ruleName
		}
		// 规则产出未过校验（不应发生）→ 继续走 LLM 兜底
	}
	if client != nil {
		for attempt := 0; attempt < 2; attempt++ {
			out, err := client.Generate(llm.SystemPrompt, norm)
			if err != nil {
				continue
			}
			out = stripCodeFence(out)
			if ok, _ := validate.Valid(out); ok {
				st.LLM++
				return out, "llm"
			}
		}
		st.LLMFallbackFixed++
		return rules.FixedAnswer, "llm-invalid"
	}
	// 无 LLM 配置且规则未命中：诚实降级为固定答案
	st.LLMFallbackFixed++
	return rules.FixedAnswer, "no-rule"
}

// stripCodeFence 剥离模型可能输出的 Markdown 代码围栏。
func stripCodeFence(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```fofa")
	s = strings.TrimPrefix(s, "```")
	if idx := strings.Index(s, "```"); idx >= 0 {
		s = s[:idx]
	}
	return strings.TrimSpace(s)
}

// VerifyIntegrity 最终验收：题数一致、题号一一对应、查询非空。
func VerifyIntegrity(questions []Question, out Template) error {
	if len(out.Answers) != len(questions) {
		return fmt.Errorf("答案条数 %d != 题目条数 %d", len(out.Answers), len(questions))
	}
	qset := make(map[string]bool, len(questions))
	for _, q := range questions {
		qset[q.ID] = true
	}
	seen := map[string]bool{}
	for _, a := range out.Answers {
		if !qset[a.ID] {
			return fmt.Errorf("未知题号: %s", a.ID)
		}
		if seen[a.ID] {
			return fmt.Errorf("重复题号: %s", a.ID)
		}
		seen[a.ID] = true
		if strings.TrimSpace(a.Query) == "" {
			return fmt.Errorf("题号 %s 查询语句为空", a.ID)
		}
	}
	if out.PkgID == "" {
		return fmt.Errorf("参赛包编号为空")
	}
	return nil
}
