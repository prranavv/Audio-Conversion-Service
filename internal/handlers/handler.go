package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
	"github.com/prranavv/peritys_submission/internal/helpers"
)

// upgrader to manage WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin for simplicity
	},
}

// Handler is a struct where all the handler methods hang off of
type Handler struct{}

// NewHandler is a constructor that makes the Handler struct
func NewHandler() *Handler {
	return &Handler{}
}

// HandleAudioStream is a handler that handles the conversion of WAV files to FLAC files
func (h *Handler) HandleAudioStream(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer ws.Close()

	// Read and validate the WAV header
	_, headerData, err := ws.ReadMessage()
	if err != nil {
		log.Println("Error reading WAV header:", err)
		return
	}
	headerReader := bytes.NewReader(headerData)
	header, err := helpers.ReadHeader(headerReader)
	if err != nil {
		log.Println("Invalid WAV header:", err)
		return
	}
	audioDataReader := &helpers.WebSocketAudioReader{
		Ws: ws,
	}
	// Define the ffmpegCmdFunc for production
	cmd := exec.Command(
		"ffmpeg",
		"-f", "s16le", // Raw PCM data
		"-ar", fmt.Sprint(header.SampleRate), // Sample rate from header
		"-ac", fmt.Sprint(header.NumChannels), // Number of channels
		"-i", "pipe:0",
		"-f", "flac",
		"pipe:1",
	)
	if err := helpers.ProcessAudioStream(ws, header, audioDataReader, cmd); err != nil {
		log.Println("Error Processing Audio:", err)
	}
}
