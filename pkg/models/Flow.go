package models

type Flow struct {
	ClientId   int64 `json:"client_id"`
	ExportFlow int64 `json:"export_flow"` // 出口流浪
	InletFlow  int64 `json:"inlet_flow"`  // 入口流量
	FlowLimit  int64 `json:"flow_limit"`  // 流量限制
}
