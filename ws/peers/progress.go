package peers

import (
	"errors"
	"io"

	"github.com/gorilla/websocket"
	"github.com/onee-only/miner-node/lib"
)

type WebSocketReader struct {
	conn *websocket.Conn
}

func (r *WebSocketReader) Read(p []byte) (int, error) {
	messageType, payload, err := r.conn.ReadMessage()
	if messageType != websocket.BinaryMessage || err != nil {
		lib.HandleErr(errors.New("wtf this should not happen"))
	}
	if len(payload) == 0 {
		return 0, io.EOF
	}
	return copy(p, payload), nil
}
