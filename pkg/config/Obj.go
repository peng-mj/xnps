package config

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	Remark     string   `json:"-" toml:"Remark"`
	InitTime   int64    `json:"init_time" toml:"InitTime"`
	BasePath   string   `json:"base_path" toml:"BasePath"`
	DbType     string   `json:"database_type" toml:"DbType"`
	AppKeys    []string `json:"app_keys" toml:"AppKeys"`
	BridgePort int      `json:"bridge_port" toml:"BridgePort"`
	WebPort    int      `json:"web_port" toml:"WebPort"`
	UsagePorts [][]int  `json:"usage_ports" toml:"UsagePorts"`
}

type UsagePort struct {
	ports [][]int
}

func NewPorts(ports [][]int) *UsagePort {
	return &UsagePort{ports: ports}
}

func (u *UsagePort) String() string {
	bd := strings.Builder{}
	ports := mergeRanges(u.ports)
	for i := 0; i < len(ports); i++ {
		s := ""
		if len(ports[i]) > 1 {
			s = fmt.Sprintf("%d-%d", ports[i][0], ports[i][len(ports[i])-1])
		} else {
			s = fmt.Sprintf("%d-%d", ports[i][0], ports[i][0])
		}
		bd.WriteString(s)
		if i != len(ports) {
			bd.WriteString(",")
		}
	}
	return bd.String()
}
func (u *UsagePort) Format() *UsagePort {
	ports := u.ports
	var res [][]int
	for i := range ports {
		if len(ports[i]) > 0 {
			res = append(res, []int{ports[i][0], ports[i][len(ports[i])-1]})
		}
	}
	res = mergeRanges(res)
	u.ports = res
	return u
}
func (u *UsagePort) Ports() [][]int {
	return u.ports
}

func (u *UsagePort) Load(ports string) error {
	ss := strings.Split(ports, ",")
	sls := make([][]int, len(ss))
	for i := range ss {
		s := strings.Split(ss[i], "-")
		if len(s) != 2 {
			return fmt.Errorf("%v is not a format string", ss[i])
		}
		start, err := strconv.Atoi(s[0])
		if err != nil {
			return fmt.Errorf("%v is not a number string", s[0])
		}
		end, err := strconv.Atoi(s[0])
		if err != nil {
			return fmt.Errorf("%v is not a number string", s[0])
		}
		if start > end {
			start, end = end, start
		}
		sls = append(sls, []int{start, end})
	}
	u.ports = mergeRanges(sls)
	return nil
}
func mergeRanges(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	var merged [][]int
	cur := intervals[0]
	for i := 0; i < len(intervals); i++ {
		if cur[1]+1 >= intervals[i][0] {
			cur[1] = intervals[i][1]
		} else {
			merged = append(merged, cur)
			cur = intervals[i]
		}
	}
	// 将最后一个范围加入结果集合
	merged = append(merged, cur)

	return merged
}
