package peers

import "sync"

type TPeers struct {
	V map[string]*Peer
	M sync.Mutex
}

var Peers *TPeers = &TPeers{
	V: make(map[string]*Peer),
	M: sync.Mutex{},
}
