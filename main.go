package main

import (
	"log"
	"p2pdistributefilestorage/p2p"
)

func main() {
	opts := p2p.TcpTransportOptions {
		ListenAddr: ":3000",
		ShakeHand: p2p.NopeHandShake,
		Decoder: p2p.DefaultDecoder{},
	}	
	tr := p2p.NewTcpTransport(opts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal("error accepting conn")
	}

	select {}
}
