package dto

import "errors"

type Response struct {
	Code  int         `json:"code"`
	ErMsg string      `json:"msg,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type User struct {
	Id        int32
	Uid       int64
	AuthLevel int32
	ExpireAt  int64
	OtpKey    string
	Valid     bool
}

type Page struct {
	Num  int32
	Size int32
}

func (p *Page) IsValid() (err error) {
	if p.Num < 0 || p.Size < 0 {
		return errors.New("page number or page size error, size or num should >= 0")
	}
	return
}
