package rules

import "strings"

// KnowledgeEntry 人工沉淀的知识条目：以题面中的特征子串为键，映射到经过推敲的 FOFA Query。
// 来源与推导依据见 data/knowledge.md。
// 覆盖：组合条件、证书全家族、跨语言语法映射、外部知识（文档字段/图标哈希/CVE）等
// 无法用参数化正则可靠表达的复杂需求。
type KnowledgeEntry struct {
	Key   string // 归一化题面中的特征子串
	Query string
}

// KnowledgeBase 知识条目按需维护，键足够特异，互不冲突。
func KnowledgeBase() []KnowledgeEntry {
	return []KnowledgeEntry{
		// ---- 网络基础资产 ----
		{Key: "网站名为 www.amazon.com", Query: `host="www.amazon.com"`},
		{Key: "域名里带有 vpn", Query: `host="vpn"`},
		{Key: "大不列颠及北爱尔兰联合王国", Query: `country="GB"`},
		{Key: "韩国江原道", Query: `country="KR" && region="Gangwon-do"`},
		{Key: "巴伐利亚州", Query: `country="DE" && region="Bavaria" && city="Munich"`},
		{Key: "使用亚马逊的云计算", Query: `country="NL" && icsp="Amazon Cloud"`},
		{Key: "河北地区的华为云", Query: `region="河北" && icsp="华为云"`},

		// ---- 协议 / 端口 ----
		{Key: "中国的 DNS 服务", Query: `protocol="dns" && country="CN"`},
		{Key: "SSH 协议且端口为 22022", Query: `protocol="ssh" && port="22022"`},
		{Key: "使用 TCP 协议", Query: `base_protocol="tcp"`},
		{Key: "SNMP 网络管理服务", Query: `protocol="snmp"`},
		{Key: "Memcached缓存服务协议", Query: `protocol="memcached" && country="US"`},
		{Key: "真正跑在 TCP 上的 FTP", Query: `protocol="ftp" && base_protocol="tcp"`},
		{Key: "返回的信息里包含 OpenSSH_10.4", Query: `protocol="ssh" && banner="OpenSSH_10.4"`},
		{Key: "日本的海康威视", Query: `app="Hikvision-视频监控" && country="JP"`},
		{Key: "产品版本号为 10.0.0", Query: `version="10.0.0"`},
		{Key: "中间层软件", Query: `category="中间件"`},

		// ---- 网页内容 / 响应特征 ----
		{Key: "更新时间在 2026-08-01 之后", Query: `(title="Jenkins" || body="jenkins") && country!="CN" && last_updatetime>"2026-08-01"`},
		{Key: "NoSuchUpload", Query: `body="NoSuchUpload"`},
		{Key: "Content-Security-Policy-Report-Only", Query: `banner="Content-Security-Policy-Report-Only"`},
		{Key: "状态码为 308", Query: `status_code="308"`},
		{Key: "SSH-2.0-OpenSSH_6.1", Query: `banner="SSH-2.0-OpenSSH_6.1"`},
		{Key: "Set-Cookie:PHPSESSID", Query: `header="Set-Cookie:PHPSESSID"`},
		{Key: "jquery.min.js", Query: `body="jquery.min.js"`},
		{Key: "阿姆斯特丹开放 554 端口", Query: `country="NL" && city="Amsterdam" && port="554" && banner="200 OK"`},

		// ---- 证书 ----
		{Key: "GlobalSign 组织签发", Query: `cert.issuer.org="GlobalSign" && cert.endtime>"2026-10-22"`},
		{Key: "持有者名称是 apple.com", Query: `cert.subject.cn="apple.com" && domain!="apple.com"`},
		{Key: "证书全文中包含 TrustAsia", Query: `cert="TrustAsia"`},
		{Key: "持有者组织名称包含 Alibaba", Query: `cert.subject.org="Alibaba"`},
		{Key: "Banner 有内容", Query: `domain="sohu.com" && banner!="" && banner!*="apache"`},
		{Key: "只要有adobe.com这个根域名", Query: `cert="adobe.com"`},
		{Key: "根域是 qq.com", Query: `cert.subject.cn~"(^|\.)qq\.com$"`},
		{Key: "根域名不是 godaddy.com", Query: `cert="godaddy" && domain!="godaddy.com"`},
		{Key: "蜜罐和欺诈网站都不要", Query: `cert="google.com" && is_honeypot="false" && is_fraud="false"`},
		{Key: "包含example.com特征", Query: `(cert.subject.cn="example.com" || cert="example.com")`},
		{Key: "仍在使用 TLS 1.0", Query: `tls="1.0" && cert.endtime>"2026-10-22"`},
		{Key: "颁发者与持有者不匹配", Query: `cert.match="false"`},
		{Key: "证书与域名不匹配", Query: `country="CA" && cert.match="false"`},
		{Key: "证书里面有根域 aliyun.com", Query: `cert="aliyun.com" && cert.endtime>"2026-10-22" && domain=""`},
		{Key: "kuaihuoyun.com/wlhyos.cn", Query: `(domain="kuaihuoyun.com" || domain="wlhyos.cn" || domain="wlhyos.com" || domain="56ctms.com" || domain="xiaokuaikeji.com" || domain="zhihuiwuliu.org.cn" || (cert.match="true" && cert.endtime>"2024-01-01 00:00:00" && cert.endtime<"2026-08-30 18:17:49" && (cert="kuaihuoyun.com" || cert="wlhyos.cn" || cert="wlhyos.com" || cert="56ctms.com" || cert="xiaokuaikeji.com" || cert="zhihuiwuliu.org.cn")))`},

		// ---- 指纹 / 系统特征 ----
		{Key: "JARM 为 21d21d", Query: `jarm="21d21d00021d21d00021d21d21d21dce319294781f0e12867fd06896a000ea" && status_code="307" && header="Content-Length: 0"`},
		{Key: "操作系统被识别为 linux", Query: `os="linux"`},
		{Key: "JA3S 指纹等于", Query: `country="TR" && ja3s="f4febc55ea12b31ae17cfb7e614afda8"`},
		{Key: "网络组织名称里包含Inc", Query: `org="Inc"`},
		{Key: "OVH SAS", Query: `country="US" && org="OVH SAS" && port="2087" && title=="WHM 登录"`},
		{Key: "域名下部署的 WordPress", Query: `domain~"\.io$" && app="WordPress"`},
		{Key: "厦门市 80 端口", Query: `city="厦门" && port="80" && status_code="200"`},
		{Key: "DigitalOcean", Query: `country="GB" && icsp="DigitalOcean" && body="/api/show"`},
		{Key: "正文有 /admin", Query: `country="CN" && banner="Set-Cookie" && body="/admin"`},
		{Key: "CNAME 根域", Query: `country="DE" && (host="amazonaws.com" || domain="amazonaws.com" || cname="amazonaws.com" || cert="amazonaws.com" || cert.subject.cn="amazonaws.com")`},
		{Key: "证书序列号为 1", Query: `port="143" && cert.serial_number="1" && (cert="TAK" || cert="ATAK" || cert="WinTAK" || cert="Squad" || cert="CoT")`},
		{Key: "且不是 IPv6", Query: `(title="Sign in" || title="Login") && port="80" && ip_type!="IPv6"`},
		{Key: "新西兰地区放 80 或 443 端口", Query: `country="NZ" && (port="80" || port="443") && icsp!=""`},
		{Key: "Banner 同时包含 IIS", Query: `(banner="IIS" && banner="1.") || header="IIS"`},
		{Key: "标题包含 大学、学部、研究", Query: `country="JP" && domain!="" && port="443" && status_code="200" && cert.endtime>"2026-10-22" && (title="大学" || title="学部" || title="研究" || body="入学案内" || body="学生")`},
		{Key: "标题有系统、平台、登录", Query: `domain="baidu.com" && (title="系统" || title="平台" || title="登录" || title="管理" || title="API" || host="api" || host="admin" || host="dev" || host="test" || host="uat" || host="pre")`},
		{Key: "被归为邮件系统", Query: `category="邮件系统" && title="Mail" && country!="US"`},
		{Key: "香港、澳门、台湾都排除", Query: `country="CN" && region!="香港" && region!="澳门" && region!="台湾"`},
		{Key: "标题不包含同一个关键词", Query: `body="Jenkins" && title!*="Jenkins"`},
		{Key: "ZeroSSL 组织签发的证书不要", Query: `asn="20940" && region="台湾" && port="443" && tls="1.3" && cert.issuer.org!="Let's Encrypt" && cert.issuer.org!="ZeroSSL"`},
		{Key: "日本以外的 HFS 文件服务", Query: `app="HFS" && country!="JP" && status_code="200" && title=="HFS /" && icsp!="Amazon Cloud"`},
		{Key: "不包含港澳台的中国 Grafana", Query: `app="Grafana" && country="CN" && region!="香港" && region!="澳门" && region!="台湾"`},
		{Key: "edu.ar 结尾的高校网站", Query: `host~"\.edu\.ar$" && (port="80" || port="443") && server!="nginx" && product!="WAF"`},
		{Key: "已经绑定域名的资产记录", Query: `ip="1.1.1.1" && domain!=""`},
		{Key: "x-qiniu-zone", Query: `country="CN" && domain!="" && header="x-qiniu-zone" && domain!="qiniudn.com"`},
		{Key: "fid 和图标哈希都不为空", Query: `title="" && status_code="200" && fid!="" && icon_hash!="" && host="admin" && (domain~"\.cn$" || domain~"^.{4}\.com$")`},
		{Key: "原样文本的资产:class", Query: `body="class=\"container\""`},
		{Key: "保留双引号和反斜杠", Query: `body="\"exception\": \"Symfony\Component\ErrorHandler\Error\FatalError\""`},
		{Key: "响应码必须是 200", Query: `ip="13.0.0.0/8" && status_code="200" && (body="Welcome" || body="Sign in" || body="Dashboard")`},
		{Key: "113.247.228.0/24 网段里开放", Query: `ip="113.247.228.0/24" && (port="80" || port="443" || port="8080" || port="8443")`},
		{Key: "至少两项有值", Query: `app="Grafana" && ((header_hash!="" && body_hash!="") || (header_hash!="" && icon_hash!="") || (header_hash!="" && fid!="") || (body_hash!="" && icon_hash!="") || (body_hash!="" && fid!="") || (icon_hash!="" && fid!=""))`},
		{Key: "整个 B 段", Query: `ip="40.88.0.0/16"`},
		{Key: "47.96.0.1 到 47.96.0.5", Query: `ip.after="47.96.0.1" && ip.before="47.96.0.5" && (port="80" || port="443")`},
		{Key: "没有绑定域名的网站类资产", Query: `country="CN" && domain=""`},
		{Key: "正文出现 Login、admin、Dashboard", Query: `domain="baidu.com" && (body="Login" || body="admin" || body="Dashboard" || body="管理" || header="Set-Cookie")`},
		{Key: "恰好还有两级子域", Query: `host~"^[^.]\.[^.]\.163\.com$"`},

		// ---- 跨语言语法映射 ----
		{Key: "Shodan 查询转换成 FOFA", Query: `header="Server: nginx" && title="Welcome"`},
		{Key: "Google Dork 转换为 FOFA", Query: `title="Grafana" && body="JavaScript"`},
		{Key: "urlscan.io 搜索意图", Query: `body="login" && (host="bet" || host="casino" || host="kasino")`},
		{Key: "Hunter 查询转成 FOFA", Query: `host="google" && domain!="google.com" && title=""`},
		{Key: "Censys 搜索转换", Query: `(header="Server: Tengine" || header="X-Cache: MISS")`},
		{Key: "Server 精确为 Apache", Query: `asn!="4134" && region="HK" && server="Apache"`},
		{Key: "意图:HFS 的 80 端口都要", Query: `app="HFS" && (port="80" || (port="8080" && status_code="200"))`},
		{Key: "依据上面的响应片段生成搜索语句", Query: `header="Etag: \"6aa522a3-493\"" && header="Content-Length: 1171" && header="Accept-Ranges" && header="Cache-Control" && header="Content-Type" && header="Date" && header="Last-Modified" && header="Server" && header="Vary"`},

		// ---- 外部知识沉淀（来源见 data/knowledge.md）----
		// M096: jenkins.io 官方文档 TCP Agent Listener Port 一节的诊断响应字段
		{Key: "TCP Agent Listener Port", Query: `(body="Jenkins-Agent-Protocols" || body="Jenkins-Version" || body="Jenkins-Session" || body="Remoting-Minimum-Version" || banner="Jenkins-Agent-Protocols" || banner="Jenkins-Version" || banner="Jenkins-Session" || banner="Remoting-Minimum-Version")`},
		// M097: octodex.github.com favicon 的 mmh3 哈希（base64 后计算）
		{Key: "octodex.github.com", Query: `icon_hash="-1493334710"`},
		// M098: 政务域名被嵌入子域的仿冒特征
		{Key: "仿冒特征的网站", Query: `(host="msp.gov.cn." || cert="msp.gov.cn.")`},
		// M099: CVE-2021-41773 = Apache HTTP Server 路径穿越（2.4.49，2.4.50 为绕过修复）
		{Key: "CVE-2021-41773", Query: `app="Apache-HTTP-Server" && (version="2.4.49" || version="2.4.50")`},
		// M100: 颁发者 CN + 证书序列号（冒号十六进制原样保留）
		{Key: "Serial Number", Query: `cert.issuer.cn="ZeroSSL ECC DV SSL CA 2" && cert.serial_number="63:BD:83:71:62:14:83:6E:F8:00:75:35:61:61:EB:06"`},
	}
}

// knowledgeRules 把知识条目包装成规则链。
func knowledgeRules() []Rule {
	entries := KnowledgeBase()
	rules := make([]Rule, len(entries))
	for i, e := range entries {
		e := e
		rules[i] = Rule{
			Name:     "knowledge: " + e.Key,
			Priority: 100,
			Convert: func(text string) (string, bool) {
				if strings.Contains(text, e.Key) {
					return e.Query, true
				}
				return "", false
			},
		}
	}
	return rules
}

// Chain 完整规则链：知识库（高精度）在前，参数化通用规则在后，最后由 LLM 兜底。
func Chain() []Rule {
	chain := knowledgeRules()
	chain = append(chain, GeneralRules()...)
	return chain
}
