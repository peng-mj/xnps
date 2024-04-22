package dto

import "net/http"

const (
	ErrParam    = 1000
	ErrNotFound = 1001
	ErrPasswd   = 1002

	NeedOtpCode = 2000
)

type RspCode int

func (m RspCode) String() string {
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
	case ErrPasswd:
		return "password or username error"

	case NeedOtpCode:
		return "need otp code"

	default:
		return "unknown type"
	}
}
func (m RspCode) Int() int {
	return int(m)
}
