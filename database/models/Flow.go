package models

import "sync"

type Flow struct {
	ClientId   int64 `json:"ClientId"`
	ExportFlow int64 `json:"ExportFlow"` //出口流浪
	InletFlow  int64 `json:"InletFlow"`  //入口流量
	FlowLimit  int64 `json:"FlowLimit"`  //流量限制
	sync.RWMutex
}

func (s *Flow) Add(in, out int64) {
	s.Lock()
	defer s.Unlock()
	s.InletFlow += in
	s.ExportFlow += out
}
