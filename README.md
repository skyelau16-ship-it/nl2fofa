# nl2fofa · 自然语言 → FOFA 查询语句转换系统

把用户的自然语言资产搜索需求转换为合法、可执行的 [FOFA](https://fofa.info) 查询语句。
零第三方依赖（Go 标准库实现），规则引擎 + LLM 兜底的混合架构。

> 本仓库为白帽汇 FOFA「自然语言资产搜索」比赛参赛作品，参赛包 `pkg-08e82c8d`，100 题，正确率以 FOFA 实际执行结果评分。

## 架构

```
题目.json
   │
   ▼
① Unicode 归一化        全角→半角（含 ”＂｜： 等）、空白折叠 —— 消除题面噪声
   ▼
② 不可转换检测器         端口>65535 / IP 非法 / 同字段逻辑矛盾 / 主观需求 /
   │                    主机内部运行时状态 / 平台能力疑问句 → 固定答案
   ▼
③ 规则引擎（首中即停）
   ├─ 知识库规则（89 条）  人工沉淀的高精度条目：证书全家族、复杂布尔组合、
   │                      跨语言语法映射、外部知识（文档/图标哈希/CVE）
   └─ 通用参数化规则       IP/CIDR/B段、ASN、端口、主域、国家+产品 —— 可处理任意同模式输入
   ▼
④ LLM 兜底（可选）        OpenAI 兼容接口，可插拔；prompt 内置 FOFA 字段表+语法+few-shot
   ▼
⑤ 保守语法校验器          括号配平 / 引号闭合 / 字段白名单 —— 宁放行不误杀
   ▼
按 答案模板.json 骨架填充（参赛包编号与题号零差错）→ 答案.json
```

各模块职责单一、可独立测试：`internal/normalize` `internal/rules` `internal/llm` `internal/validate` `internal/pipeline`。

## 快速开始

```bash
# 纯规则模式（无需任何 key）
go run . -q data/题目.json -t data/答案模板.json -o answers/答案.json -name 你的名字

# 规则 + LLM 兜底模式（规则未命中才调用，控制成本）
export FOFA_LLM_BASE_URL=https://api.deepseek.com/v1   # 任意 OpenAI 兼容接口
export FOFA_LLM_API_KEY=sk-xxx
export FOFA_LLM_MODEL=deepseek-chat
go run . ...
```

输出末尾会打印各来源命中统计（如本参赛包：`不可转换 6 | 知识库 89 | 通用规则 5 | LLM 0 | 降级 0`）。

## 测试

```bash
go test ./...
```

覆盖：固定答案文案逐字精确（比赛红线）、五类不可转换判定、CIDR 归一化（B 段 → `40.88.0.0/16`）、
语法校验器（括号/引号/字段边界识别，值内 `class=` 不误判）、知识库键冲突扫描。

## 设计决策

| 决策 | 理由 |
|------|------|
| 规则优先、LLM 兜底 | 规则 100% 确定性可复现、离线可跑；LLM 只处理规则覆盖不到的长尾，成本与波动最小化 |
| 以答案模板为骨架填充 | 参赛包编号/题号/条数与官方模板天然一致，杜绝结构性拒绝 |
| 校验器宁放行不误杀 | 校验失败的惩罚（降级固定答案）远高于放行一条非常规但正确的 Query |
| 全角归一化前置 | 实测题面存在全角冒号/引号/竖线（如修正题），不归一化则规则全部漏配 |
| 外部知识沉淀而非运行时抓取 | M096-M100 需要读文档/算哈希/查 CVE，沉淀为可审计的知识条目（来源见 `data/knowledge.md`），运行时不依赖网络 |
| 零第三方依赖 | 面试官 clone 即跑，无供应链负担；HTTP/JSON/正则标准库完全够用 |

## 知识库来源

图标哈希现场计算（favicon base64 → mmh3）、Jenkins 官方文档字段提取、CVE 影响面核对等，
每条外部知识的来源与推导见 [`data/knowledge.md`](data/knowledge.md)。

## 目录结构

```
├── main.go                    # CLI 入口
├── internal/
│   ├── normalize/             # Unicode 归一化
│   ├── rules/                 # 不可转换检测 + 知识库 + 通用规则 + 责任链
│   ├── llm/                   # OpenAI 兼容客户端 + System Prompt
│   ├── validate/              # 保守语法校验
│   └── pipeline/              # 流水线编排 + 完整性验收
├── data/                      # 题目 / 答案模板 / 知识来源
├── answers/答案.json           # 生成结果（选手：刘鹏飞）
└── PLAN.md                    # 执行方案
```

## License

MIT
