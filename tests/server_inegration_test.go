package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prranavv/peritys_submission/internal/helpers"
	"github.com/prranavv/peritys_submission/internal/server"
)

func TestAudioStreamWebSocketHandler(t *testing.T) {
	serverAddr := "localhost:8080"
	server := server.Run()
	defer server.Shutdown(context.Background())
	// Give the server a moment to start
	time.Sleep(1 * time.Second)

	u := url.URL{Scheme: "ws", Host: serverAddr, Path: "/stream"}
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v,%v", err, u.String())
	}
	defer ws.Close()

	header, audioData, err := prepareTestWAVData()
	if err != nil {
		t.Fatalf("Failed to prepare test WAV data: %v", err)
	}
	err = ws.WriteMessage(websocket.BinaryMessage, header)
	if err != nil {
		t.Fatalf("Failed to send WAV header: %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		buf := make([]byte, 4096)
		for {
			n, err := audioData.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Errorf("Error reading test audio data: %v", err)
				return
			}
			if n > 0 {
				err := ws.WriteMessage(websocket.BinaryMessage, buf[:n])
				if err != nil {
					t.Errorf("Error sending audio data: %v", err)
					return
				}
			}
		}
	}()

	var receivedData bytes.Buffer
	go func() {
		defer wg.Done()
		for {
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					break
				}
				t.Errorf("Error reading from WebSocket: %v", err)
				return
			}

			if messageType == websocket.BinaryMessage {
				receivedData.Write(message)
			}
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	if receivedData.Len() == 0 {
		t.Fatal("No data received from server")
	}

	err = os.WriteFile("test_output.flac", receivedData.Bytes(), 0644)
	if err != nil {
		t.Fatalf("Failed to write output file: %v", err)
	}

	// Use ffprobe to verify the output file
	cmd := exec.Command("ffprobe", "test_output.flac")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ffprobe failed: %v, output: %s", err, string(output))
	}
	// Clean up the output file
	os.Remove("test_output.flac")

}

func prepareTestWAVData() ([]byte, io.Reader, error) {
	// Create a simple WAV header
	header := &helpers.WAVHeader{
		// Fill in appropriate fields
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     36 + 8, // Adjust as needed
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16,
		AudioFormat:   1,    // PCM
		NumChannels:   1,    // Mono
		SampleRate:    8000, // 8 kHz
		ByteRate:      8000 * 2,
		BlockAlign:    2,
		BitsPerSample: 16,
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: 16000, // 2 seconds of audio
	}

	// Serialize the header
	headerBuf := new(bytes.Buffer)
	err := binary.Write(headerBuf, binary.LittleEndian, header)
	if err != nil {
		return nil, nil, err
	}

	// Generate dummy audio data (e.g., a sine wave or silence)
	audioData := make([]byte, header.Subchunk2Size)
	// Fill audioData with sample values (e.g., zeros for silence)

	// Return header bytes and audio data reader
	return headerBuf.Bytes(), bytes.NewReader(audioData), nil
}
