package helpers

import (
	"bytes"
	"io"
	"log"
)

type Command interface {
	Start() error
	Wait() error
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
}

type MockCommand struct {
	StdinWriter  *io.PipeWriter
	StdoutReader *io.PipeReader
}

func NewMockCommand() *MockCommand {
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	mockCmd := &MockCommand{
		StdinWriter:  stdinWriter,
		StdoutReader: stdoutReader,
	}

	// Simulate processing: copy data from stdin to stdout, possibly modifying it
	go func() {
		defer stdoutWriter.Close()
		defer stdinReader.Close()
		buf := make([]byte, 4096)
		for {
			n, err := stdinReader.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("Error in reading:", err)
			}
			// Simulate processing (e.g., convert to uppercase)
			processedData := bytes.ToUpper(buf[:n])
			stdoutWriter.Write(processedData)
		}
	}()

	return mockCmd
}

func (m *MockCommand) Start() error {
	// No action needed for the mock
	return nil
}

func (m *MockCommand) Wait() error {
	// No action needed for the mock
	return nil
}

func (m *MockCommand) StdinPipe() (io.WriteCloser, error) {
	return m.StdinWriter, nil
}

func (m *MockCommand) StdoutPipe() (io.ReadCloser, error) {
	return m.StdoutReader, nil
}
