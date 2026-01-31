package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *Rpc) error
}

type GOBDecoder struct{}

func (g GOBDecoder) Decode(r io.Reader, msg *Rpc) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, msg *Rpc) error {
	buf := make([]byte, 1024)
	n,err := r.Read(buf)	
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]

	return nil
}
