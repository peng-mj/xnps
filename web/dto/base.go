package dto

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

type Page struct {
	Page int32
	Size int32
}

func (p *Page) Check() (err error) {

	return
}
