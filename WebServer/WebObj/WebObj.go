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

// 用于控制平台，不同用户的最大的同时登录使用数量
type Salt struct {
	SaltMap map[string]string
	maxLen  int
}

func NewSaltMap(maxLen int) *Salt {
	salt := Salt{
		SaltMap: make(map[string]string, maxLen),
		maxLen:  maxLen}
	return &salt
}

func (s *Salt) Get(user string) (string, error) {
	if s.SaltMap != nil {
		if salt, ok := s.SaltMap[user]; ok {
			return salt, nil
		}
	}
	return "", errors.New("user " + user + "have no salt")
}
func (s *Salt) Add(user, salt string) error {
	if len(s.SaltMap) < s.maxLen {
		log.Println("map数量：", len(s.SaltMap))
		s.SaltMap[user] = salt
		return nil
	}
	return errors.New("user " + user + "have no salt")

}
func (s *Salt) Del(user string) {
	delete(s.SaltMap, user)
}
