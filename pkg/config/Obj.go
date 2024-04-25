package config

type Config struct {
	Remark     string   `json:"-" toml:"Remark"`
	InitTime   int64    `json:"init_time" toml:"InitTime"`
	BasePath   string   `json:"base_path" toml:"BasePath"`
	DbType     string   `json:"database_type" toml:"DbType"`
	AppKeys    []string `json:"AppKeys" toml:"AppKeys"`
	WebPort    int      `json:"WebPort" toml:"WebPort"`
	UsagePorts [][]int  `json:"UsagePorts" toml:"UsagePorts"`
}
