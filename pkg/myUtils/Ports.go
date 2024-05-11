package myUitls

import (
	"fmt"
	"strconv"
	"strings"
	"tunpx/pkg/crypt"
)

type Ports struct {
	ports [][]int
}

func NewPorts(ports [][]int) *Ports {
	return &Ports{ports: ports}
}

func (u *Ports) String() string {
	bd := strings.Builder{}
	ports := crypt.MergeRanges(u.ports)
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
func (u *Ports) Format() *Ports {
	ports := u.ports
	var res [][]int
	for i := range ports {
		if len(ports[i]) > 0 {
			res = append(res, []int{ports[i][0], ports[i][len(ports[i])-1]})
		}
	}
	res = crypt.MergeRanges(res)
	u.ports = res
	return u
}
func (u *Ports) Ports() [][]int {
	return u.ports
}

func (u *Ports) Load(ports string) error {
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
	u.ports = crypt.MergeRanges(sls)
	return nil
}
