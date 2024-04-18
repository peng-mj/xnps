package server

import (
	"net"
	"sync"
)

type Connection struct {
	Conn       net.Conn
	Close      chan struct{}
	ConnType   uint32
	MaxConnNum uint32
	CurConnNum uint32
	Status     uint32
}

func NewConnection() {

}

func (conn Connection) CLose() {
	close(conn.Close)
}
func (conn Connection) Run() {

}
func (conn Connection) Restart() {

}
func (conn Connection) AddClient() {

}
func (conn Connection) RemoveClient() {

}

type Server struct {
	ConnectMap  map[int]*Connection
	ConnectPool []Connection
	RwLock      sync.RWMutex
	//mp          sync.Map
}
