package helpers

import (
	"bytes"
	"testing"
)

func TestProcessAudioStream(t *testing.T) {
	// Prepare mock WebSocket
	mockWS := &MockWebSocket{}

	// Prepare WAV header
	header := &WAVHeader{
		SampleRate:    44100,
		NumChannels:   2,
		Subchunk2Size: 8, // Small size for testing
	}

	// Prepare audio data
	audioData := bytes.NewReader([]byte("testdata"))

	// Create a MockCommand
	mockCmd := NewMockCommand()

	// Call the function under test
	err := ProcessAudioStream(mockWS, header, audioData, mockCmd)
	if err != nil {
		t.Fatalf("ProcessAudioStream returned error: %v", err)
	}

	// Verify that messages were sent over the WebSocket
	if len(mockWS.Messages) == 0 {
		t.Fatal("No messages were sent over the WebSocket")
	}

	// Verify the content of the messages
	expectedOutput := "TESTDATA" // Since we convert input to uppercase
	receivedOutput := string(bytes.Join(mockWS.Messages, nil))
	if receivedOutput != expectedOutput {
		t.Fatalf("Expected output %q, got %q", expectedOutput, receivedOutput)
	}
}
