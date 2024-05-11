package dto

import "errors"

type GroupReq struct {
	Name       string  `json:"name,omitempty"`
	UsagePorts [][]int `json:"usage_ports,omitempty"`
	MaxClient  int32   `json:"max_client,omitempty"`
	Remark     string  `json:"remark,omitempty"`
}

type GroupGetReq struct {
	Filters struct {
		Ids  []string `json:"ids"`
		Name string   `json:"name"`
		Uid  int64    `json:"uid"`
	} `json:"filters,omitempty"`
	Page Page `json:"page,omitempty"`
}

func (g *GroupReq) IsValid() (err error) {
	if len(g.Name) < 3 {
		return errors.New("name length max > 3")
	}
	return nil
}
