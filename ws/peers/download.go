package peers

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"os"

	"github.com/gorilla/websocket"
	"github.com/onee-only/miner-node/config"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/ws/messages"
	"github.com/schollz/progressbar/v3"
)

const chunkSize = 1024

func StartDownloadingBlockChain() {
	config.IsDownloading = true
	peer := getRandomPeer()

	m := messages.Message{
		Kind:    messages.MessageDownloadRequest,
		Payload: nil,
	}

	bytes, err := json.Marshal(m)
	lib.HandleErr(err)
	peer.Inbox <- bytes
}

func downloadBlockchain(header []byte, conn *websocket.Conn) {
	fileSize := int64(binary.LittleEndian.Uint64(header))

	bar := progressbar.DefaultBytes(
		fileSize,
		"Downloading blockchain",
	)

	_, err := io.Copy(io.MultiWriter(file, bar), &WebSocketReader{conn})
	lib.HandleErr(err)

	bar.Finish()
	config.IsDownloading = false
}

func uploadBlockchain(p *Peer) {
	config.IsDownloading = true
	file, err := os.Open("blockchain.db")
	lib.HandleErr(err)
	defer file.Close()

	fileInfo, err := file.Stat()
	lib.HandleErr(err)
	fileSize := make([]byte, 8)
	binary.LittleEndian.PutUint64(fileSize, uint64(fileInfo.Size()))
	err = p.Conn.WriteMessage(websocket.BinaryMessage, fileSize)
	lib.HandleErr(err)

	for {
		buffer := make([]byte, chunkSize)
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				lib.HandleErr(err)
			}
			break
		}
		err = p.Conn.WriteMessage(websocket.BinaryMessage, buffer[:bytesRead])
		lib.HandleErr(err)
	}
	err = p.Conn.WriteMessage(websocket.BinaryMessage, nil)
	lib.HandleErr(err)
	config.IsDownloading = false
}
