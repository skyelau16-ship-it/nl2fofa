// Package validate 对生成的 FOFA Query 做保守语法校验。
//
// 策略：宁放行、不误杀 —— 校验失败意味着题目会被降级为固定答案，
// 因此只检查确定性的错误（括号/引号失衡、未知字段名），
// 正则、转义等拿不准的情况一律放行。
package validate

import (
	"regexp"
	"strings"
)

// knownFields FOFA 官方字段白名单（含本项目知识库用到的全部字段）。
var knownFields = map[string]bool{
	"ip": true, "ip.after": true, "ip.before": true,
	"port": true, "port_size": true, "domain": true, "host": true, "cname": true,
	"title": true, "body": true, "banner": true, "header": true, "header_hash": true, "body_hash": true,
	"status_code": true, "protocol": true, "base_protocol": true,
	"country": true, "country_code": true, "region": true, "city": true,
	"asn": true, "org": true, "icsp": true,
	"icon_hash": true, "fid": true, "jarm": true, "ja3": true, "ja3s": true,
	"product": true, "app": true, "category": true, "version": true,
	"os": true, "server": true, "tls": true, "last_updatetime": true, "ip_type": true,
	"cert": true, "cert.subject.cn": true, "cert.subject.org": true,
	"cert.issuer.cn": true, "cert.issuer.org": true,
	"cert.starttime": true, "cert.endtime": true,
	"cert.serial_number": true, "cert.match": true,
	"is_honeypot": true, "is_fraud": true,
}

// fieldRe 只在表达式边界（行首/括号/空白/逻辑运算符）后识别字段名，
// 避免把引号值里的内容（如 body="class=..." 中的 class=）误判为字段。
var fieldRe = regexp.MustCompile(`(?:\A|[(\s]|\|\||&&)\s*([A-Za-z_][A-Za-z0-9_.]*)\s*(==|!=|\*=|!\*=|!~|~|>=|<=|>|<|=)`)

// unescapedQuotes 统计未被 \ 转义的引号数量（Go RE2 不支持 lookbehind，手工扫描）。
func unescapedQuotes(q string) int {
	n, escaped := 0, false
	for _, r := range q {
		if escaped {
			escaped = false
			continue
		}
		switch r {
		case '\\':
			escaped = true
		case '"':
			n++
		}
	}
	return n
}

// Valid 校验一条 FOFA Query 的语法，返回是否合法与失败原因。
func Valid(query string) (bool, string) {
	q := strings.TrimSpace(query)
	if q == "" {
		return false, "空查询"
	}
	// 括号配平
	depth := 0
	for _, r := range q {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
			if depth < 0 {
				return false, "括号失衡"
			}
		}
	}
	if depth != 0 {
		return false, "括号失衡"
	}
	// 引号成对（忽略转义）
	if n := unescapedQuotes(q); n%2 != 0 {
		return false, "引号不成对"
	}
	// 字段名必须在白名单内
	for _, m := range fieldRe.FindAllStringSubmatch(q, -1) {
		if !knownFields[m[1]] {
			return false, "未知字段: " + m[1]
		}
	}
	return true, ""
}
