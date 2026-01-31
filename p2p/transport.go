
package p2p 
// Peer represent the remote node
type Peer interface {
	Close() error
}


// Transport anything that handle the communication in 
// node
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan Rpc
} 