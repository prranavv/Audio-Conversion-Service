package helpers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// WavHeaderSize is the size of the WAV header
const WavHeaderSize = 44

// WAVHeader is a struct that represents the 44 byte WAV header
type WAVHeader struct {
	ChunkID       [4]byte // "RIFF"
	ChunkSize     uint32
	Format        [4]byte // "WAVE"
	Subchunk1ID   [4]byte // "fmt "
	Subchunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Subchunk2ID   [4]byte // "data"
	Subchunk2Size uint32
}

// ReadHeader is a method that reads the header of the WAV file and validates it.
func ReadHeader(r io.Reader) (*WAVHeader, error) {
	header := &WAVHeader{}
	headerBytes := make([]byte, WavHeaderSize)
	totalRead := 0

	// Read the header in a loop to handle streaming data
	for totalRead < WavHeaderSize {
		n, err := r.Read(headerBytes[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += n
	}

	buffer := bytes.NewReader(headerBytes)
	if err := binary.Read(buffer, binary.LittleEndian, header); err != nil {
		return nil, err
	}

	// Validate header
	if string(header.ChunkID[:]) != "RIFF" || string(header.Format[:]) != "WAVE" {
		return nil, errors.New("invalid WAV file format")
	}

	return header, nil
}
