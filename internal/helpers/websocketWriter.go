package helpers

// Define an interface for WebSocket communication
type WebSocketWriter interface {
	WriteMessage(messageType int, data []byte) error
}

type MockWebSocket struct {
	Messages     [][]byte
	MessageTypes []int
	WriteErr     error
}

func (m *MockWebSocket) WriteMessage(messageType int, data []byte) error {
	if m.WriteErr != nil {
		return m.WriteErr
	}
	m.MessageTypes = append(m.MessageTypes, messageType)
	m.Messages = append(m.Messages, append([]byte(nil), data...)) // Copy data
	return nil
}
