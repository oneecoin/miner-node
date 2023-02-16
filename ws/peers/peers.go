package peers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/gorilla/websocket"
	"github.com/onee-only/miner-node/config"
	"github.com/onee-only/miner-node/lib"
)

type TPeers struct {
	V map[string]*Peer
	M sync.Mutex
	C chan config.ChanMessageType
}

var Peers *TPeers = &TPeers{
	V: make(map[string]*Peer),
	M: sync.Mutex{},
	C: make(chan config.ChanMessageType),
}

func Connect() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "Connecting to nodes "
	s.FinalMSG = "Nodes connected!\n"
	s.Start()

	var peerList []string
	res, err := http.Get(fmt.Sprintf("http://%s/peers?port=%d", config.MempoolAddress, config.Port))
	lib.HandleErr(err)
	err = json.NewDecoder(res.Body).Decode(&peerList)
	lib.HandleErr(err)

	for _, address := range peerList {
		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws?port=%d&publicKey=%s", address, config.Port, config.PublicKey), nil)
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
	s.Stop()
}

func (*TPeers) InitPeer(p *Peer) {
	Peers.M.Lock()
	defer Peers.M.Unlock()
	go p.read()
	go p.write()
	Peers.V[p.GetAddress()] = p
}
