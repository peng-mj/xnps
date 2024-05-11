package dto

import "net/http"

const (
	ErrParam       = 1000
	ErrNotFound    = 1001
	ErrPasswd      = 1002
	ErrRateLimit   = 1003
	ErrDbError     = 1004
	ErrInternalErr = 1005

	ErrCreateConfigFile = 1100
	ErrCreateUser       = 1101
	ErrCreateGroup      = 1102

	NeedOtpCode = 2000
)

type RspCode int

func (m RspCode) String() string {
	res := http.StatusText(int(m))
	if res != "" {
		return res
	}
	switch m {
	// 连接参数
	case ErrParam:
		return "the input parameter is incorrect"
	case ErrNotFound:
		return "data not found"
	case ErrPasswd:
		return "password or username error"
	case ErrRateLimit:
		return "rate limit"
	case ErrDbError:
		return "database error"
	case ErrInternalErr:
		return "internal error"

	case NeedOtpCode:
		return "need otp code"

	default:
		return "unknown type"
	}
}
func (m RspCode) Int() int {
	return int(m)
}
