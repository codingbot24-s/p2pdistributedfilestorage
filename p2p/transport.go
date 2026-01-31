
package p2p 
// Peer represent the remote node
type Peer interface {}


// Transport anything that handle the communication in 
// node
type Transport interface {
	ListenAndAccept() error
} 