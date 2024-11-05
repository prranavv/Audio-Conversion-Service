package helpers

import (
	"io"
	"log"

	"github.com/gorilla/websocket"
)

// ProcessAudioStream processes audio data from an io.Reader, converts it using ffmpeg,
// and sends the converted data back over a WebSocket connection.
func ProcessAudioStream(
	ws WebSocketWriter,
	header *WAVHeader,
	audioData io.Reader,
	cmd Command, // Changed from ffmpegCmdFunc to cmd Command
) error {
	// Set up pipes
	ffmpegStdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println("Error creating stdin pipe:", err)
		return err
	}

	ffmpegStdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Error creating stdout pipe:", err)
		return err
	}

	expectedDataSize := int(header.Subchunk2Size)
	totalBytesRead := 0

	// Start the command
	if err := cmd.Start(); err != nil {
		log.Println("Error starting command:", err)
		return err
	}

	// Goroutine to write audio data to command's stdin
	writeErrCh := make(chan error, 1)
	go func() {
		defer ffmpegStdin.Close()
		buf := make([]byte, 4096)
		for {
			if expectedDataSize > 0 && totalBytesRead >= expectedDataSize {
				log.Println("Expected audio data received.")
				break
			}

			n, err := audioData.Read(buf)
			if err != nil {
				if err == io.EOF {
					log.Println("Finished reading audio data.")
					break
				}
				log.Println("Error reading audio data:", err)
				writeErrCh <- err
				return
			}

			if n > 0 {
				totalBytesRead += n
				if _, err := ffmpegStdin.Write(buf[:n]); err != nil {
					log.Println("Error writing to command stdin:", err)
					writeErrCh <- err
					return
				}
			}
		}
		writeErrCh <- nil
	}()

	// Reading loop from command's stdout
	buf := make([]byte, 4096)
	for {
		n, err := ffmpegStdout.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("Command stdout closed")
				break
			}
			log.Println("Error reading from command stdout:", err)
			return err
		}
		if n > 0 {
			if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				log.Println("Error writing to WebSocket:", err)
				return err
			}
		}
	}

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		log.Println("Command process exited with error:", err)
		return err
	}

	// Check for errors from the goroutine
	if writeErr := <-writeErrCh; writeErr != nil {
		return writeErr
	}

	// Close the WebSocket connection
	log.Println("Closing WebSocket connection")
	return nil
}
