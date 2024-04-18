package web

import "xnps/pkg/cache"

type MiddleBase struct {
	secret *cache.Cache
	salt   *cache.Cache
}
type AuthUser struct {
	MiddleBase
	User     string
	UserId   int32
	AuthCode int32
	IsAdmin  bool
}

func (m *MiddleBase) AuthMiddle() {

}

func (m *MiddleBase) JwtMiddle() {

}
