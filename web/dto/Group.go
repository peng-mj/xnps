package dto

import "errors"

type GroupCreateRequest struct {
	Name       string  `json:"name,omitempty"`
	UsagePorts [][]int `json:"usage_ports,omitempty"`
	MaxClient  int32   `json:"max_client,omitempty"`
	Remark     string  `json:"remark,omitempty"`
}

func (g *GroupCreateRequest) IsValid() (err error) {
	if len(g.Name) < 3 {
		return errors.New("name length max > 3")
	}

	return nil
}
