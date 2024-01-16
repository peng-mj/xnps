package SysConfig

type SysConfig struct {
	driver *Driver
	SystemConfig
}

// 只执行一次，当数据库中无配置信息时
type SystemConfig struct {
	WebHost           string `json:"webHost"` //默认服务地址
	WebPort           int64  `json:"webPort"` //对外服务默认8912
	WebUserName       string `json:"webUserName"`
	WebPassword       string `json:"webPassword"`    //因为是sha256加密，所以需要考虑密码重置的情况
	WebOpenCaptcha    bool   `json:"webOpenCaptcha"` //是否开启验证码校验
	AuthCryptKey      string `json:"authCryptKey"`
	AllowPorts        string `json:"allowPorts"`
	PublicKey         string `json:"publicKey"`
	BridgeType        string `json:"bridgeType"` //tcp、udp、socket、kcp
	BridgePort        int    `json:"bridgePort"`
	BridgeHost        string `json:"bridgeHost"` //
	LogLevel          int    `json:"logLevel"`
	LogPath           string `json:"logPath"`
	MaxClient         int    `json:"maxClient"` //这里可以根据不同性能设备做一下说明
	MaxConn           int    `json:"maxConn"`
	DisConnTimeoutSec int    `json:"disConnTimeoutSec"`
	AllowRegistration bool   `json:"allowRegistration"`
}

// 只执行一次，当数据库中无配置信息时
type ClientConfig struct {
	WebHost           string `json:"webHost"` //默认服务地址
	WebPort           int64  `json:"webPort"` //对外服务默认8912
	WebUserName       string `json:"webUserName"`
	WebPassword       string `json:"webPassword"`    //因为是sha256加密，所以需要考虑密码重置的情况
	WebOpenCaptcha    bool   `json:"webOpenCaptcha"` //是否开启验证码校验
	AuthCryptKey      string `json:"authCryptKey"`
	AllowPorts        string `json:"allowPorts"`
	PublicKey         string `json:"publicKey"`
	BridgeType        string `json:"bridgeType"` //tcp、udp、socket、kcp
	BridgePort        int    `json:"bridgePort"`
	BridgeHost        string `json:"bridgeHost"` //
	LogLevel          int    `json:"logLevel"`
	LogPath           string `json:"logPath"`
	MaxClient         int    `json:"maxClient"` //这里可以根据不同性能设备做一下说明
	MaxConn           int    `json:"maxConn"`
	DisConnTimeoutSec int    `json:"disConnTimeoutSec"`
	AllowRegistration bool   `json:"allowRegistration"`
}
