package netTool

import (
	"errors"
	"fmt"
	"net"
)

// CheckPortOccupied to check net port
func CheckPortOccupied(port int, protocol string) error {
	switch protocol {
	case "tcp":
	case "udp":
	default:
		return errors.New("not support " + protocol)
	}
	_, err := net.Listen(protocol, fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return nil
}
