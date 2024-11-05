package helpers

import (
	"bytes"
	"io"
	"log"

	"github.com/gorilla/websocket"
)

// WebSocketAudioReader implements the io.Reader interface
type WebSocketAudioReader struct {
	Ws     *websocket.Conn //Ws is a Websocket conection
	buffer bytes.Buffer
}

// Read function is a method that satisfies the io.Reader interface
func (r *WebSocketAudioReader) Read(p []byte) (int, error) {
	if r.buffer.Len() > 0 {
		return r.buffer.Read(p)
	}
	messageType, data, err := r.Ws.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			log.Println("WebSocket closed unexpectedly:", err)
			return 0, io.EOF
		}
		return 0, err
	}

	if messageType != websocket.BinaryMessage {
		return 0, nil
	}

	if _, err := r.buffer.Write(data); err != nil {
		return 0, err
	}
	return r.buffer.Read(p)
}
