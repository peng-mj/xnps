package dto

import (
	"errors"
	"strings"
	"tunpx/pkg/config"
	"tunpx/pkg/crypt"
	"tunpx/pkg/sysTool"
)

type ConfigReq struct {
	OrgName      string  `json:"org_name,omitempty"`
	BasePath     string  `json:"base_path"`
	WebPort      int     `json:"web_port"`
	BridgePort   int     `json:"bridge_port"`
	UsagePorts   [][]int `json:"usage_ports"`
	Level        int32   ` json:"level,omitempty"`
	Username     string  `json:"username"`
	Password     string  `json:"password,omitempty"`
	Emile        string  `json:"emile,omitempty"`
	EmileAuthKey string  `json:"emile_auth_key"` // SMTP or other
	OTAKeys      string  ` json:"ota_keys,omitempty"`
	ExpirationAt int64   ` json:"expiration_at"`
	MaxConn      int     ` json:"max_conn"`
}

func (c *ConfigReq) Validity() error {
	if !crypt.CheckEmail(c.Username) {
		return errors.New("username must be an email address")
	}

	if c.WebPort > 65535 || c.WebPort < 80 {
		return errors.New("web port error. maybe 8900 ok")
	}
	if c.BridgePort > 65535 || c.BridgePort < 80 {
		return errors.New("bridge port error. maybe 8901 ok")
	}
	if crypt.CheckPassed(c.Password) < 3 {
		return errors.New("password to week, include at least three types of numbers, uppercase letters, lowercase letters, and special symbols")
	}
	usagePort := config.NewPorts(c.UsagePorts).Format()
	if len(usagePort.Ports()) == 0 {
		return errors.New("usage proxy port should not be null, check it now")
	}
	c.UsagePorts = usagePort.Ports()

	if c.MaxConn < 3 {
		c.MaxConn = 3
	}
	if c.BasePath == "" {
		c.BasePath = "./data"
	}

	c.BasePath = strings.TrimRight(c.BasePath, "/")
	sysTool.CreateFolder(c.BasePath)
	return nil

}
