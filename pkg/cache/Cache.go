package cache

import (
	"sync"
)

type KvMap struct {
	mp sync.Map
}

func NewMp() *KvMap {
	return &KvMap{
		mp: sync.Map{},
	}
}

func (s *KvMap) Get(k string) (any, bool) {
	return s.mp.Load(k)
}
func (s *KvMap) Put(k string, v any) {
	s.mp.Store(k, v)
}
func (s *KvMap) Del(k string) {
	s.mp.Delete(k)
}
