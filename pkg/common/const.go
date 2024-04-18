package common

const (
	CONN_DATA_SEQ     = "*#*" //Separator
	VERIFY_EER        = "vkey"
	VERIFY_SUCCESS    = "sucs"
	WORK_MAIN         = "main" //主连接
	WORK_CHAN         = "chan" //隧道
	WORK_CONFIG       = "conf" //通过配置文件连接
	WORK_REGISTER     = "rgst" //注册
	WORK_SECRET       = "sert" //私密代理
	WORK_FILE         = "file" //文件代理，后续去掉
	WORK_P2P          = "p2pm" //p2p，后续去掉
	WORK_P2P_VISITOR  = "p2pv"
	WORK_P2P_PROVIDER = "p2pp"
	WORK_P2P_CONNECT  = "p2pc"
	WORK_P2P_SUCCESS  = "p2ps"
	WORK_P2P_END      = "p2pe"
	WORK_STATUS       = "stus"
	RES_CLOSE         = "clse"
	NEW_UDP_CONN      = "udpc" //p2p udp conn
	NEW_TASK          = "task"
	NEW_CONF          = "conf"
	CONN_TCP          = "tcp"
	CONN_UDP          = "udp"
	CONN_TEST         = "TST"

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

type MsgType uint16

func (m *MsgType) GetMsg() string {
	switch *m {
	//连接参数
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
	case 100:
		return "tcp"
	case 101:
		return "udp"
	case 102:
		return "kcp"
	case 103:
		return "p2p"
	//连接类型（加密）
	case 200:
		return "s-tcp"
	case 201:
		return "s-udp"
	case 202:
		return "s-kcp"
	case 203:
		return "s-p2p"
	case 1000:
		return "success"
	case 2000:
		return "client error"
	case 3000:
		return "server error"

	default:
		return "unknown type"
	}
}

func IsTunnelMode(md string) bool {
	return md == MODE_P2P || md == MODE_TCP || md == MODE_UDP || md == MODE_HTTP || md == MODE_HTTPS || md == MODE_SECRECT
}
