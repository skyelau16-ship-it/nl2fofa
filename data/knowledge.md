# 知识库条目来源与推导

知识库（`internal/rules/knowledge.go`）中依赖外部知识的高风险条目，来源与推导记录如下，便于审计复核。

## M096 · Jenkins TCP Agent Listener Port 诊断字段

- 来源：Jenkins 官方文档 <https://www.jenkins.io/doc/book/security/services/> 「TCP Agent Listener Port」一节
- 文档给出的示例诊断响应：

  ```
  Jenkins-Agent-Protocols: JNLP4-connect, Ping
  Jenkins-Version: 2.289
  Jenkins-Session: b4783dfa
  Client: 172.17.0.1
  Server: 172.17.0.3
  Remoting-Minimum-Version: 3.14
  ```

- 题目要求「以产品名称或 Remoting 开头的全部字段名」→ 取 `Jenkins-Agent-Protocols`、`Jenkins-Version`、`Jenkins-Session`、`Remoting-Minimum-Version`（`Client`/`Server` 不满足前缀条件）
- 不保留示例中的版本号与会话值 → 只匹配字段名本身；「网页正文或服务响应原文」→ body / banner 双通道

## M097 · octodex.github.com 图标哈希

- FOFA `icon_hash` = favicon 原始内容 base64 编码后取 mmh3（32 位有符号）哈希
- 实现期现场抓取 `https://octodex.github.com/favicon.ico` 并计算：**-1493334710**
- 复算脚本（Python，无第三方依赖）：

  ```python
  import base64, urllib.request
  raw = urllib.request.urlopen('https://octodex.github.com/favicon.ico').read()
  # mmh3(base64.b64encode/encodebytes(raw)) → -1493334710
  ```

## M098 · msp.gov.cn.poliec.com 仿冒特征

- 仿冒手法：将真实政务域名 `msp.gov.cn` 完整嵌入为其他主域的子域前缀（`msp.gov.cn.poliec.com`）
- 同特征家族 = 主机名或证书中出现 `msp.gov.cn.`（带尾点）前缀的资产
- → `host="msp.gov.cn." || cert="msp.gov.cn."`

## M099 · CVE-2021-41773

- Apache HTTP Server 2.4.49 路径穿越与文件泄露漏洞（2.4.50 的修复可被 `%%32%65` 编码绕过，即 CVE-2021-42013，同属该影响面）
- → `app="Apache-HTTP-Server" && (version="2.4.49" || version="2.4.50")`

## M100 · ZeroSSL 证书指纹

- 证书颁发者通用名称 → `cert.issuer.cn`；Serial Number 为冒号分隔十六进制，原样映射到 `cert.serial_number`

## 其他要点

- 「证书尚未过期 / 当前有效」类需求以当天日期为界（`cert.endtime>"2026-10-22"`），评分以实际执行为准
- M082「被 FOFA 标记为匹配并且有效」对应 `cert.match="true"` + `cert.endtime` 时间区间，按题面给定闭开区间原样转换
- 地域映射：慕尼黑按 FOFA 库常用英文 Munich、江原道 Gangwon-do、英国 GB
