package p2p

import (
	"net"
	"sync"
)

type TcpTransport struct {
	ListenAddr string
	Listener   net.Listener

	mu			sync.Mutex
	Peers		map[net.Addr]Peer	 		
}

func NewTcpTransport (lsaddr string) *TcpTransport {
	return &TcpTransport{
		ListenAddr:  lsaddr,
	}
}