package WebObj

import (
	"errors"
	"log"
)

const (
	CmdClientAdd    = "ClientAdd"
	CmdClientEdit   = "ClientEdit"
	CmdClientDelete = "ClientDel"
	CmdTunnelAdd    = "TunnelAdd"
	CmdTunnelEdit   = "TunnelEdit"
	CmdTunnelDel    = "TunnelDel"
	CmdGroupAdd     = "GroupAdd"
	CmdGroupEdit    = "GroupEdit"
	CmdGroupDel     = "GroupDel"

	ErrNoData       = "NoData"
	ErrOverLoginNum = "OverLoginNumber"

	//Err = 0
	//ErrNoData = 0
	//ErrNoData = 0
	//ErrNoData = 0
	//ErrNoData = 0
	//ErrNoData = 0
)

type Request struct {
	MsgType string      `json:"Type,omitempty"`
	Data    interface{} `json:"Data,omitempty"`
}
type LoginData struct {
}

type DoLogin struct {
	Username string `json:"Username"`
	//Timestamp int64  `json:"Timestamp"`
}
type Login struct {
	Username string `json:"Username,omitempty"`
	Password string `json:"Password,omitempty"` //sha256加密
}

//type SysConfig struct {
//	Username string `json:"Username,omitempty"`
//	Password string `json:"Password,omitempty"` //sha256加密
//}

// 用于控制平台，不同用户的最大的同时登录使用数量
type KVManage struct {
	SaltMap map[string]string
	maxLen  int
}

func NewKVMap(maxLen int) *KVManage {
	salt := KVManage{
		SaltMap: make(map[string]string, maxLen),
		maxLen:  maxLen}
	return &salt
}

func (s *KVManage) Get(user string) (string, error) {
	if s.SaltMap != nil {
		if salt, ok := s.SaltMap[user]; ok {
			return salt, nil
		}
	}
	return "", errors.New("user " + user + "have no salt")
}
func (s *KVManage) Add(user, salt string) error {
	if len(s.SaltMap) < s.maxLen {
		log.Println("map数量：", len(s.SaltMap))
		s.SaltMap[user] = salt
		return nil
	}
	return errors.New("user " + user + "have no salt")

}
func (s *KVManage) Del(user string) {
	delete(s.SaltMap, user)
}

//func NewTokenManege(delaySec int64, maxTokenNum int) *TokenManager {
//
//	tokenManege := TokenManager{
//		Tokens:         make(chan Token, 20),
//		Ticker:         time.NewTicker(time.Duration(delaySec) * time.Second),
//		ExpirationTime: delaySec,
//	}
//	return &tokenManege
//}

//
//type Token struct {
//	Key string
//	AddTime int64
//}
//
//type TokenManager struct {
//	KvManage KVManage
//	Tokens chan Token
//	Ticker *time.Ticker
//	ExpirationTime int64
//	MaxToken int
//}
//
//func (t *TokenManager) AutoDelToken() {
//
//	select {
//	case <-t.Ticker.C:
//
//		t.Ticker.Reset()
//
//	}
//
//}
//
//func (t *TokenManager) AddToken(token string) error {
//	if t.Tokens.Get(token)
//}
//
//func (t *TokenManager) DelToken(token string) {
//
//}
