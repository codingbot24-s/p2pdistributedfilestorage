package p2p


type HandShaker func(peer Peer) error

func NopeHandShake(peer Peer) error { return nil }
