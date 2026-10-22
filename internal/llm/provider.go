// Package llm 提供 OpenAI 兼容接口的可插拔客户端（DeepSeek/Qwen/GLM/Kimi 等均可）。
// 零第三方依赖，直接用标准库实现 /chat/completions 调用。
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client OpenAI 兼容聊天客户端。
type Client struct {
	BaseURL string // 如 https://api.deepseek.com/v1
	APIKey  string
	Model   string
	HTTP    *http.Client
}

// FromEnv 从环境变量构建客户端；未配置 key 时返回 nil（规则-only 模式）。
// FOFA_LLM_BASE_URL / FOFA_LLM_API_KEY / FOFA_LLM_MODEL
func FromEnv() *Client {
	key := os.Getenv("FOFA_LLM_API_KEY")
	if key == "" {
		return nil
	}
	base := os.Getenv("FOFA_LLM_BASE_URL")
	if base == "" {
		base = "https://api.deepseek.com/v1"
	}
	model := os.Getenv("FOFA_LLM_MODEL")
	if model == "" {
		model = "deepseek-chat"
	}
	return &Client{BaseURL: base, APIKey: key, Model: model,
		HTTP: &http.Client{Timeout: 60 * time.Second}}
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Generate 调用模型生成单条结果，temperature=0 保证可复现。
func (c *Client) Generate(system, user string) (string, error) {
	body, _ := json.Marshal(chatRequest{
		Model: c.Model,
		Messages: []message{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		Temperature: 0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimRight(c.BaseURL, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.Error != nil {
		return "", fmt.Errorf("llm api: %s", out.Error.Message)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("llm api: 空响应")
	}
	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}
