// nl2fofa 自然语言 → FOFA 查询语句转换系统（规则引擎 + LLM 兜底）。
//
// 用法：
//	go run . -q data/题目.json -t data/答案模板.json -o answers/答案.json -name 选手名称
//
// LLM 配置（可选，未配置则纯规则模式）：
//	export FOFA_LLM_BASE_URL=https://api.deepseek.com/v1
//	export FOFA_LLM_API_KEY=sk-xxx
//	export FOFA_LLM_MODEL=deepseek-chat
//	或写入项目根目录 config.env（已 gitignore）。
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/skyelau16-ship-it/nl2fofa/internal/llm"
	"github.com/skyelau16-ship-it/nl2fofa/internal/pipeline"
)

func main() {
	questionsPath := flag.String("q", "data/题目.json", "题目 JSON 路径")
	tplPath := flag.String("t", "data/答案模板.json", "答案模板 JSON 路径")
	outPath := flag.String("o", "answers/答案.json", "输出答案 JSON 路径")
	player := flag.String("name", "", "选手名称")
	flag.Parse()

	loadDotEnv("config.env")

	var questions []pipeline.Question
	mustLoad(*questionsPath, &questions)
	var tpl pipeline.Template
	mustLoad(*tplPath, &tpl)

	client := llm.FromEnv()
	mode := "纯规则模式（未配置 LLM key）"
	if client != nil {
		mode = "规则 + LLM 兜底模式（" + client.Model + "）"
	}
	fmt.Println("运行模式:", mode)

	out, stats := pipeline.Run(questions, tpl, *player, client)

	if err := pipeline.VerifyIntegrity(questions, out); err != nil {
		fmt.Fprintln(os.Stderr, "完整性校验失败:", err)
		os.Exit(1)
	}

	f, err := os.Create(*outPath)
	must(err)
	w := bufio.NewWriter(f)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false) // 中文/引号原样输出
	enc.SetIndent("", "  ")
	must(enc.Encode(out))
	must(w.Flush())
	must(f.Close())

	fmt.Printf("完成：%d 题 | 不可转换 %d | 知识库 %d | 通用规则 %d | LLM %d | 降级固定答案 %d\n",
		len(out.Answers), stats.Impossible, stats.Knowledge, stats.General, stats.LLM, stats.LLMFallbackFixed)
	fmt.Println("已写入:", *outPath)
}

func mustLoad(path string, v any) {
	b, err := os.ReadFile(path)
	must(err)
	must(json.Unmarshal(b, v))
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// loadDotEnv 读取 KEY=VALUE 格式的本地配置（不存在则跳过），密钥绝不入库。
func loadDotEnv(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if k, v, ok := strings.Cut(line, "="); ok {
			k, v = strings.TrimSpace(k), strings.TrimSpace(v)
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}
