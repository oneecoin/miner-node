package peers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
)

type TPeers struct {
	V map[string]*Peer
	M sync.Mutex
}

var Peers *TPeers = &TPeers{
	V: make(map[string]*Peer),
	M: sync.Mutex{},
}

func Connect() {
	s := lib.CreateSpinner(
		"Connecting to nodes",
		"Nodes connected!",
	)

	var peerList []string
	res, err := http.Get(fmt.Sprintf("http://%s/peers?port=%d", properties.MempoolAddress, properties.Port))
	lib.HandleErr(err)
	err = json.NewDecoder(res.Body).Decode(&peerList)
	lib.HandleErr(err)

	for _, address := range peerList {
		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws?port=%d&publicKey=%s", address, properties.Port, properties.PublicKey), nil)
		lib.HandleErr(err)

		p := &Peer{
			Conn:  conn,
			Inbox: make(chan []byte),
		}

		sep := strings.Split(address, ":")
		p.Address.Host = sep[0]
		p.Address.Port = sep[1]

		Peers.InitPeer(p)
	}
	go listenBlockBroadcast()
	s.Stop()
}

func (*TPeers) InitPeer(p *Peer) {
	Peers.M.Lock()
	defer Peers.M.Unlock()
	go p.read()
	go p.write()
	Peers.V[p.GetAddress()] = p
}
