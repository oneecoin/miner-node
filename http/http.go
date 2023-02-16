package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/onee-only/miner-node/config"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/ws/peers"
)

var wsUpgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
}

var prs *peers.TPeers = peers.Peers

func InitServer(port int) {

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		publicKey := r.URL.Query().Get("publicKey")
		port := r.URL.Query().Get("port")
		host := strings.Split(r.RemoteAddr, ":")[0]
		address := fmt.Sprintf("%s:%s", host, port)

		wsUpgrader.CheckOrigin = func(r *http.Request) bool {
			// send http request to the address
			res, err := http.Get("http://" + address + "/check")
			if err != nil {
				return false
			}
			if res.StatusCode != http.StatusAccepted {
				return false
			}

			defer res.Body.Close()
			a := struct{ PublicKey string }{}
			err = json.NewDecoder(res.Body).Decode(&a)
			if err != nil {
				return false
			}
			if a.PublicKey != publicKey {
				return false
			}
			return true
		}
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		lib.HandleErr(err)

		p := &peers.Peer{
			Conn:  conn,
			Inbox: make(chan []byte),
			Address: peers.TAddress{
				Host: host,
				Port: port,
			},
		}

		prs.InitPeer(p)
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			bytes, err := json.Marshal(struct{ PublicKey string }{
				PublicKey: config.PublicKey,
			})
			lib.HandleErr(err)
			w.WriteHeader(http.StatusAccepted)
			w.Write(bytes)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	lib.HandleErr(err)
}
