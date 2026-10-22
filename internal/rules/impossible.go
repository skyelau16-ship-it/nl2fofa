package rules

import (
	"regexp"
	"strconv"
	"strings"
)

// Impossible 检测「不能转换为 FOFA 搜索语句」的需求（README 明确要求填固定答案）。
// 覆盖五类：非法端口值、非法 IP、同字段逻辑矛盾、主观评价、主机内部运行时状态。
func Impossible(text string) bool {
	return invalidPort(text) || invalidIP(text) || contradicted(text) ||
		subjective(text) || hostInternalState(text) || metaQuestion(text)
}

// invalidPort 端口号超出 0-65535（README 示例：70001 / 本包 M091：1000000）。
func invalidPort(text string) bool {
	for _, m := range regexp.MustCompile(`(\d{1,7})\s*端口`).FindAllStringSubmatch(text, -1) {
		if n, err := strconv.Atoi(m[1]); err == nil && n > 65535 {
			return true
		}
	}
	return false
}

var ipLikeRe = regexp.MustCompile(`\b(\d{1,3}(?:\.\d{1,3}){3,5})(?:/\d{1,2})?\b`)

// invalidIP 形如 1.2.3.4.5 的地址/网段（M092）：IPv4 至多 4 段，每段 0-255。
func invalidIP(text string) bool {
	for _, m := range ipLikeRe.FindAllStringSubmatch(text, -1) {
		parts := strings.Split(m[1], ".")
		if len(parts) > 4 {
			return true
		}
		for _, p := range parts {
			if n, err := strconv.Atoi(p); err != nil || n > 255 {
				return true
			}
		}
	}
	return false
}

var contradictionRe = regexp.MustCompile(`必须等于\s*([A-Za-z0-9]+)[^。]*又必须不等于\s*([A-Za-z0-9]+)`)

// contradicted 同一约束既要求等于又要求不等于（M093：US 且非 US）。
// RE2 不支持反向引用，用捕获组手动比对。
func contradicted(text string) bool {
	for _, m := range contradictionRe.FindAllStringSubmatch(text, -1) {
		if m[1] == m[2] {
			return true
		}
	}
	return false
}

// subjective 主观评价类需求（M094「最安全的网站」）——无客观判定标准。
func subjective(text string) bool {
	return strings.Contains(text, "最安全") || strings.Contains(text, "最好") ||
		strings.Contains(text, "最危险")
}

// hostInternalState 主机内部运行时状态（M095 容器内存使用率）——
// 网络空间测绘只采集对外暴露的服务特征，拿不到主机内监控指标。
func hostInternalState(text string) bool {
	return strings.Contains(text, "内存使用率") || strings.Contains(text, "CPU使用率") ||
		strings.Contains(text, "容器监控")
}

// metaQuestion 针对 FOFA 平台能力的疑问句（M015「fofa有识别…吗？」）——
// 是提问而非资产搜索需求。归一化后问号可能是半角或全角。
func metaQuestion(text string) bool {
	t := strings.ToLower(text)
	if !(strings.HasPrefix(t, "fofa有") || strings.HasPrefix(t, "fofa 是否")) {
		return false
	}
	return strings.HasSuffix(t, "吗？") || strings.HasSuffix(t, "吗?")
}
