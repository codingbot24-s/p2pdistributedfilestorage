package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTcpTransport(t *testing.T) {
	listenaddr := "4000"
	tr := NewTcpTransport(listenaddr)
	assert.Equal(t, tr.ListenAddr, listenaddr)

	tr.ListenAndAccept()
}

