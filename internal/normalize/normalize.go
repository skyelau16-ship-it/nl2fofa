// Package normalize 对自然语言输入做统一的 Unicode 归一化，
// 消除全角字符、中文标点（如 M089 原查询里的 ”4134″ 与全角 ｜）带来的匹配噪声。
package normalize

import "strings"

// fullwidth 将全角 ASCII（！＂＃..～）与常见中文标点转为半角。
func fullwidth(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r == 0x3000: // 全角空格
			r = ' '
		case r >= 0xFF01 && r <= 0xFF5E: // ！~ ～
			r -= 0xFEE0
		case r == '“' || r == '”' || r == '‘' || r == '’':
			r = '"'
		case r == '—' || r == '–':
			r = '-'
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Text 返回归一化后的文本：全角转半角、连续空白折叠为单个空格、去首尾空白。
func Text(s string) string {
	s = fullwidth(s)
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}
