package peers

import (
	"io"

	"github.com/gorilla/websocket"
	"github.com/schollz/progressbar/v3"
)

type progressWriter struct {
	writer io.Writer
	bar    *progressbar.ProgressBar
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	if err != nil {
		return n, err
	}

	pw.bar.Add(n)

	return n, nil
}

type WebSocketReader struct {
	conn *websocket.Conn
}

func (r *WebSocketReader) Read(p []byte) (int, error) {
	_, reader, err := r.conn.NextReader()
	if err != nil {
		return 0, err
	}
	return reader.Read(p)
}

type WebSocketWriter struct {
	conn *websocket.Conn
}

func (w *WebSocketWriter) Write(p []byte) (int, error) {
	err := w.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
