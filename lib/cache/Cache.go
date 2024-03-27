package cache

import (
	"errors"
	"fmt"
	"sync"
)

type KVManage struct {
	ValueMap map[string]string
	Keys     []string
	maxLen   int
	lk       sync.Map
}

func NewKVMap(maxLen int) *KVManage {
	salt := KVManage{
		ValueMap: make(map[string]string, maxLen),
		maxLen:   maxLen,
		Keys:     []string{},
	}
	return &salt
}

func (s *KVManage) Get(Key string) (string, error) {
	if s.ValueMap != nil {
		if salt, ok := s.ValueMap[Key]; ok {
			return salt, nil
		}
	}
	return "", errors.New(Key + "have no value")
}
func (s *KVManage) GetString(Key string) (string, error) {
	if s.ValueMap != nil {
		if salt, ok := s.ValueMap[Key]; ok {
			return salt, nil
		}
	}
	return "", fmt.Errorf("key=%s have no value", Key)
}
func (s *KVManage) GetWithDefault(Key, def string) string {
	if s.ValueMap != nil {
		if salt, ok := s.ValueMap[Key]; ok {
			return salt
		}
	}
	return def
}
func (s *KVManage) Put(key, value string) {
	if len(s.ValueMap) >= s.maxLen {
		s.DelLast()
	}
	s.ValueMap[key] = value
	s.Keys = append(s.Keys, key)
	return

}
func (s *KVManage) Del(Key string) {
	delete(s.ValueMap, Key)
	idx := -1
	for i, key := range s.Keys {
		if key == Key {
			idx = i
			break
		}
	}
	if idx != -1 {
		s.Keys = append(s.Keys[:idx], s.Keys[idx+1:]...)
	}
}

func (s *KVManage) DelLast() {
	if len(s.Keys) > 0 {
		delete(s.ValueMap, s.Keys[0])
		s.Keys = s.Keys[1:]
	}
}
