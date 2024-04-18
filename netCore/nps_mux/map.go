package nps_mux

import (
	"sync"
)

type ConnectMap struct {
	cMap map[int32]*conn
	sync.RWMutex
}

func NewConnMap() *ConnectMap {
	cMap := &ConnectMap{
		cMap: make(map[int32]*conn),
	}
	return cMap
}

func (s *ConnectMap) Size() (n int) {
	s.RLock()
	n = len(s.cMap)
	s.RUnlock()
	return
}

func (s *ConnectMap) Get(id int32) (*conn, bool) {
	s.RLock()
	v, ok := s.cMap[id]
	s.RUnlock()
	if ok && v != nil {
		return v, true
	}
	return nil, false
}

func (s *ConnectMap) Set(id int32, v *conn) {
	s.Lock()
	s.cMap[id] = v
	s.Unlock()
}

func (s *ConnectMap) Close() {
	for _, v := range s.cMap {
		_ = v.Close() // close all the connections in the mux
	}
}

func (s *ConnectMap) Delete(id int32) {
	s.Lock()
	delete(s.cMap, id)
	s.Unlock()
}
