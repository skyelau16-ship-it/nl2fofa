package rules

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

// GeneralRules 参数化通用规则：不依赖具体题目，可处理任意同模式输入。
// 顺序即优先级：精确匹配在前，模糊匹配在后。
func GeneralRules() []Rule {
	return []Rule{
		{Name: "ip-exact", Priority: 10, Convert: convertIPExact},
		{Name: "ip-cidr", Priority: 20, Convert: convertCIDR},
		{Name: "asn", Priority: 30, Convert: convertASN},
		{Name: "port-only", Priority: 40, Convert: convertPortOnly},
		{Name: "main-domain", Priority: 50, Convert: convertMainDomain},
		{Name: "geo-country-product", Priority: 60, Convert: convertGeoProduct},
	}
}

var (
	ipExactRe = regexp.MustCompile(`IP\s*地址[为是]?\s*(\d{1,3}(?:\.\d{1,3}){3})`)
	cidrRe    = regexp.MustCompile(`(\d{1,3}(?:\.\d{1,3}){3})/(\d{1,2})\s*网段`)
	asnRe     = regexp.MustCompile(`自治系统号[为是]?\s*(\d+)`)
	portOnly  = regexp.MustCompile(`^(?:请查询|搜索|查一下)?开放\s*(\d{1,5})\s*端口的资产`)
	mainDomRe = regexp.MustCompile(`主域[是为是]?\s*([a-z0-9-]+(?:\.[a-z0-9-]+)+)`)
)

func convertIPExact(text string) (string, bool) {
	if m := ipExactRe.FindStringSubmatch(text); m != nil {
		return fmt.Sprintf("ip=%q", m[1]), true
	}
	return "", false
}

// convertCIDR 将「x.x.x.x/NN 网段」归一化为网络地址形式的 CIDR（如 B 段 → /16）。
func convertCIDR(text string) (string, bool) {
	m := cidrRe.FindStringSubmatch(text)
	if m == nil && strings.Contains(text, "B 段") {
		if m2 := regexp.MustCompile(`(\d{1,3}(?:\.\d{1,3}){3})`).FindStringSubmatch(text); m2 != nil {
			return fmt.Sprintf("ip=%q", normalizeCIDR(m2[1]+"/16")), true
		}
	}
	if m == nil {
		return "", false
	}
	return fmt.Sprintf("ip=%q", normalizeCIDR(m[1]+"/"+m[2])), true
}

// normalizeCIDR 计算 CIDR 的网络地址（40.88.14.254/16 → 40.88.0.0/16）。
func normalizeCIDR(cidr string) string {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return cidr // 解析失败原样返回，交给校验器兜底
	}
	ones, _ := ipnet.Mask.Size()
	return ipnet.IP.String() + "/" + fmt.Sprint(ones)
}

func convertASN(text string) (string, bool) {
	if m := asnRe.FindStringSubmatch(text); m != nil {
		return fmt.Sprintf("asn=%q", m[1]), true
	}
	return "", false
}

func convertPortOnly(text string) (string, bool) {
	if m := portOnly.FindStringSubmatch(text); m != nil {
		return fmt.Sprintf("port=%q", m[1]), true
	}
	return "", false
}

func convertMainDomain(text string) (string, bool) {
	if m := mainDomRe.FindStringSubmatch(text); m != nil {
		return fmt.Sprintf("domain=%q", m[1]), true
	}
	return "", false
}

var (
	countryZh = map[string]string{
		"中国": "CN", "美国": "US", "日本": "JP", "德国": "DE", "荷兰": "NL",
		"英国": "GB", "韩国": "KR", "土耳其": "TR", "加拿大": "CA", "新西兰": "NZ",
	}
	geoProductRe = regexp.MustCompile(`(中国|美国|日本|德国|荷兰|英国|韩国|土耳其|加拿大|新西兰)(?:地区|的)?\s*(?:的)?([A-Za-z][A-Za-z0-9 -]{0,30}?)资产`)
)

// convertGeoProduct 「X国/地区的 Y 资产」→ country="X" && app~="Y"。
// 作为通用能力的演示；复杂地域组合由知识库与 LLM 兜底处理。
func convertGeoProduct(text string) (string, bool) {
	m := geoProductRe.FindStringSubmatch(text)
	if m == nil {
		return "", false
	}
	code, ok := countryZh[m[1]]
	if !ok {
		return "", false
	}
	product := strings.TrimSpace(m[2])
	if product == "" {
		return "", false
	}
	return fmt.Sprintf("country=%q && app=%q", code, product), true
}
