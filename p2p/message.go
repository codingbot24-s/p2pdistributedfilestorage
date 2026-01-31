package p2p

import "net"

// Message holds any data that is sent between two node/peers
type Message struct {
	From    net.Addr
	Payload []byte
}
