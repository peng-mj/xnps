package common

import "net/http"

const (
	CONN_DATA_SEQ = "*#*" // Separator

	UnauthorizedBytes = `HTTP/1.1 401 Unauthorized
Content-Type: text/plain; charset=utf-8
WWW-Authenticate: Basic realm="easyProxy"

401 Unauthorized`
	ConnectionFailBytes = `HTTP/1.1 404 Not Found

`
	MODE_TCP     = "tcp"
	MODE_UDP     = "udp"
	MODE_HTTP    = "http"
	MODE_HTTPS   = "https"
	MODE_SECRECT = "secrect"
	MODE_P2P     = "p2p"
)

type CodeEnum uint16

func (m CodeEnum) String() string {
	res := http.StatusText(int(m))
	if len(res) != 0 {
		return res
	}
	switch m {
	// 连接参数
	case 0:
		return "close"
	case 1:
		return "conn"
	case 2:
		return "uid"
	case 3:
		return "key"
	case 4:
		return ""
	case 5:
		return ""
	case 6:
		return ""
	case 7:
		return ""
	case 8:
		return ""
	case 9:
		return ""
	//	连接类型(正常)
	case 10:
		return "tcp"
	case 11:
		return "udp"
	case 12:
		return "kcp"
	case 13:
		return "p2p"
	// 连接类型（加密）
	case 20:
		return "s-tcp"
	case 21:
		return "s-udp"
	case 22:
		return "s-kcp"
	case 23:
		return "s-p2p"
	//	连接状态
	case 30:
		return "success"
	case 40:
		return "client error"
	case 50:
		return "server error"

	default:
		return "unknown type"
	}
}
