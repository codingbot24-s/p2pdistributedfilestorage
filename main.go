package main

import (
	"fmt"
	"log"
	"p2pdistributefilestorage/p2p"
)

func ONPeer (peer p2p.Peer) error {
	return nil
}

func main() {
	opts := p2p.TcpTransportOptions{
		ListenAddr: ":3000",
		ShakeHand:  p2p.NopeHandShake,
		Decoder:    p2p.DefaultDecoder{},
	}
	tr := p2p.NewTcpTransport(opts)

	go func() {
		for rpc := range tr.Consume() {
			fmt.Printf("received rpc %+v\n", rpc)
		}

	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal("error accepting conn")
	}

	select {}
}
