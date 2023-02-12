package socket

import (
	"encoding/json"
	"net/http"

	"github.com/onee-only/miner-node/config"
	"github.com/onee-only/miner-node/lib"
)

type publicKeyResponse struct {
	PublicKey string `json:"publicKey"`
}

func InitServer() {

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			bytes, err := json.Marshal(publicKeyResponse{PublicKey: config.PublicKey})
			lib.HandleErr(err)
			w.WriteHeader(http.StatusAccepted)
			w.Write(bytes)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":4000", nil)
}
