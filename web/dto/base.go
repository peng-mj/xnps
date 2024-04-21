package dto

import (
	"net/http"
)

const (
	ErrParam    = 1000
	ErrNotFound = 1001
)

type ErEnum int

func (m ErEnum) String() string {
	res := http.StatusText(int(m))
	if len(res) != 0 {
		return res
	}
	switch m {
	//连接参数
	case ErrParam:
		return "the input parameter is incorrect"
	case ErrNotFound:
		return "data not found"

	default:
		return "unknown type"
	}
}

type Response struct {
	Code  int         `json:"code"`
	ErMsg string      `json:"msg,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type User struct {
	Id        int32
	Uid       string
	AuthLevel int32
	ExpireAt  int64
	OtpKey    string
	Valid     bool
}

// TODO:需要参考如何定义
type Page struct {
	Offset int32
	Count  int32
}

func (p *Page) Check() (err error) {
	if p.Count > 1000 {
		p.Count = 1000
	}
	if p.Count < 0 {
		p.Count = 0
	}
	return
}
