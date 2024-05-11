package config

type Config struct {
	Remark     string   `json:"-" toml:"Remark"`
	InitTime   int64    `json:"init_time" toml:"InitTime"`
	BasePath   string   `json:"base_path" toml:"BasePath"`
	DbType     string   `json:"database_type" toml:"DbType"`
	AppKeys    []string `json:"app_keys" toml:"AppKeys"`
	BridgePort int      `json:"bridge_port" toml:"BridgePort"`
	WebPort    int      `json:"web_port" toml:"WebPort"`
	UsagePorts [][]int  `json:"usage_ports" toml:"UsagePorts"`
}
