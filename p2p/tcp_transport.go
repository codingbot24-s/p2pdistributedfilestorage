package p2p

import (
	"fmt"
	"net"
	"sync"
)


// TcpPeer represent the remote node over a tcp server
type TcpPeer struct {
	// underlying conn of peer
	conn net.Conn
	// if we dial then we are a client --> true;
	// if we accept listen then we are a server --> false
	outBond bool
}

func NewTcpPeer(conn net.Conn, outBond bool) *TcpPeer {
	return &TcpPeer{
		conn:     conn,
		outBond:  outBond,
	}
}

type TcpTransportOptions struct {
	ListenAddr string
	// shake the hand with conn 
	ShakeHand HandShaker
	Decoder Decoder	

}



type TcpTransport struct {
	TcpTransportOptions 
	Listener   net.Listener
	rwmutex  	sync.RWMutex
	Peers		map[net.Addr]Peer	 		
}
// NewTcpTransport return the new tcp transport
func NewTcpTransport (tcpopts TcpTransportOptions) *TcpTransport {
	return &TcpTransport{
		TcpTransportOptions: tcpopts,	
	}
}

// ListenAndAccept listen and accept the tcp connection
func (t *TcpTransport) ListenAndAccept () error {
	ln,err := net.Listen("tcp",t.ListenAddr)	
	if err != nil {
		return err
	}

	t.Listener = ln

	go t.startAcceptLoop()

	return nil  
}

// startAcceptLoop start the accept loop
func (t *TcpTransport) startAcceptLoop() {
	for {
		conn,err := t.Listener.Accept()
		if err != nil {
			fmt.Printf("accept error %s\n",err)
		}

		go t.handleCoon(conn)
	}
}



// hadleConn handle the connection
func (t *TcpTransport) handleCoon(conn net.Conn)  {
	fmt.Printf("new incoming connection %+v\n",conn)
	peer := NewTcpPeer(conn,true)
	if err := t.ShakeHand(peer); err != nil {
		// peer.conn.Close()
		conn.Close()
		fmt.Printf("error invalid handshake %s\n",err)
		return 
	}
	
	msg := &Rpc{}
	// start reading 
	for {
		if err := t.Decoder.Decode(conn,msg); err != nil {
			fmt.Printf("TCP decoding error %s\n",err)
			continue
		}
		msg.From = peer.conn.RemoteAddr()
		fmt.Printf("received message %+v\n",msg)
	}
}