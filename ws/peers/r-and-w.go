package peers

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/gorilla/websocket"
	"github.com/onee-only/miner-node/lib"
	"github.com/onee-only/miner-node/properties"
	"github.com/onee-only/miner-node/ws/messages"
	"github.com/schollz/progressbar/v3"
)

const chunkSize = 1024

func StartDownloadingBlockChain() {
	properties.IsDownloading = true
	peer := getRandomPeer()
	if peer == nil {
		fmt.Println(properties.WarningStr("Peer not found. Initializing DB"))
		if _, err := os.Stat("blockchain.db"); os.IsNotExist(err) {
			// file does not exist, create it
			_, err := os.Create("blockchain.db")
			lib.HandleErr(err)
		}
		return
	}

	fmt.Println(properties.WarningStr("Your DB needs to be synchronized!"))

	m := messages.Message{
		Kind:    messages.MessageDownloadRequest,
		Payload: nil,
	}

	bytes := lib.ToJSON(m)
	peer.Inbox <- bytes
	<-properties.C
}

func downloadBlockchain(header []byte, conn *websocket.Conn) {
	fileSize := int64(binary.LittleEndian.Uint64(header))

	bar := progressbar.DefaultBytes(
		fileSize,
		"Downloading blockchain",
	)

	file, err := os.Create(properties.DBName)
	lib.HandleErr(err)

	_, err = io.Copy(io.MultiWriter(file, bar), &WebSocketReader{conn})
	lib.HandleErr(err)

	bar.Finish()
	properties.C <- properties.MessageBlockchainDownloaded
}

func uploadBlockchain(p *Peer) {
	properties.C <- properties.MessageBlockchainUploading
	properties.IsDownloading = true
	file, err := os.Open(properties.DBName)
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
	properties.IsDownloading = false
}
