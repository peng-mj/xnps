package dto

type DoLogin struct {
	Username string `json:"username"`
}
type Login struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"` //sha256加密
}
