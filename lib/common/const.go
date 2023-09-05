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
	WORK_P2P_LAST     = "p2pl"
	WORK_STATUS       = "stus"
	RES_MSG           = "msg0"
	RES_CLOSE         = "clse"
	NEW_UDP_CONN      = "udpc" //p2p udp conn
	NEW_TASK          = "task"
	NEW_CONF          = "conf"
	NEW_HOST          = "host"
	CONN_TCP          = "tcp"
	CONN_UDP          = "udp"
	CONN_TEST         = "TST"
	UnauthorizedBytes = `HTTP/1.1 401 Unauthorized
Content-Type: text/plain; charset=utf-8
WWW-Authenticate: Basic realm="easyProxy"

401 Unauthorized`
	ConnectionFailBytes = `HTTP/1.1 404 Not Found

`
)
