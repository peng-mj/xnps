package dto

type LoginReq struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"` //sha256加密
	OtpCode  string `json:"otp_code,omitempty"`
}

type LoginRsp struct {
	Response
}
type UserInfo struct {
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Emile        string `json:"emile,omitempty"`
	EmileAuthKey string `json:"emile_auth_key,omitempty"`
	EmileEnable  bool   `json:"emile_enable,omitempty"`
	//MaxConn      int    `json:"maxConn,omitempty"` //maybe next feature
}
