package dto

type LoginReq struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"` //sha256加密
	OtpCode  string `json:"otp_code,omitempty"`
}

type LoginRsp struct {
	Response
}
