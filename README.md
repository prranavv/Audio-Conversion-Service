<h1 align="center">Audio Conversion Service</h1>
<p align="center"><strong>Real-Time WAV to FLAC Conversion over WebSockets</strong></p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go" />
  <img src="https://img.shields.io/badge/WebSockets-✓-blue?style=flat-square" />
  <img src="https://img.shields.io/badge/FFmpeg-✓-007808?style=flat-square" />
  <img src="https://img.shields.io/badge/Docker-✓-2496ED?style=flat-square&logo=docker" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" />
</p>

<p align="center">
  A backend microservice built in Go that receives WAV audio files over a WebSocket connection, converts them to FLAC using FFmpeg, and streams the result back to the client in real time. Dockerized, with integration and unit tests.
</p>

---

## The Problem

Audio format conversion is a common backend task, but most implementations rely on uploading a file via HTTP, waiting for server-side processing, and then downloading the result. This request-response pattern adds latency and doesn't scale well for larger files — the client is blocked until the entire conversion is complete.

## The Solution

This service uses **WebSockets** instead of HTTP for bidirectional streaming. The client sends a WAV file over the socket, the server pipes it through FFmpeg for conversion, and the resulting FLAC file is streamed back over the same connection. No file uploads, no polling, no waiting for a download link.

---

## How It Works

```
┌──────────────────────────────────────────────────────────────┐
│                         Client                               │
│                                                              │
│   ┌─────────────────────────────────────────────────────┐    │
│   │  Sends WAV file over WebSocket                      │    │
│   │  Receives FLAC file over the same connection        │    │
│   └────────────────────────┬────────────────────────────┘    │
└────────────────────────────┼─────────────────────────────────┘
                             │  WebSocket (ws://localhost:8080)
┌────────────────────────────┴─────────────────────────────────┐
│                    Go Server (port 8080)                      │
│                                                              │
│  ┌──────────┐  ┌───────────────────┐  ┌──────────────────┐   │
│  │ WebSocket│  │ ProcessAudioStream│  │  FFmpeg Process  │   │
│  │ Handler  │  │                   │  │                  │   │
│  │          │──►  Receives WAV     │──►  stdin: WAV      │   │
│  │ /stream  │  │  Pipes to FFmpeg  │  │  stdout: FLAC    │   │
│  │          │◄──  Sends FLAC back  │◄──  real-time pipe  │   │
│  └──────────┘  └───────────────────┘  └──────────────────┘   │
└──────────────────────────────────────────────────────────────┘
```

---

## Features

- **WebSocket streaming** — bidirectional audio transfer over a single persistent connection, no file upload/download round-trips
- **FFmpeg-powered conversion** — leverages FFmpeg's battle-tested audio codec support for WAV to FLAC encoding
- **Dockerized** — single `docker compose up --build` to run the entire service with all dependencies
- **Integration tests** — end-to-end test that sends a WAV file through the WebSocket and verifies the FLAC output
- **Unit tests** — `ProcessAudioStream` tested in isolation with mocked WebSocket connections and mocked FFmpeg execution
- **Manual test client** — included Node.js client script for quick manual verification

---

## API Endpoint

| Endpoint | Protocol | Description |
|---|---|---|
| **`/stream`** | WebSocket | Receives a WAV file over the connection, converts it to FLAC via FFmpeg, and sends the result back over the same connection. |

---

## Quick Start

### Prerequisites

- Docker (recommended) — or Go 1.21+ and FFmpeg installed locally

### Option 1: Docker (recommended)

```bash
git clone https://github.com/prranavv/Audio-Conversion-Servic.git
cd Audio-Conversion-Servic

docker compose up --build
```

The service starts on `http://localhost:8080`.

### Option 2: Run locally

```bash
bin/app
```

Requires FFmpeg to be installed and available on your `PATH`.

---

## Testing

### Integration test

End-to-end test that spins up the server, sends a WAV file through the WebSocket, and verifies the converted FLAC output:

```bash
go test ./tests/... -v
```

### Unit tests

Tests the `ProcessAudioStream` function in isolation with mocked WebSocket connections and mocked FFmpeg execution:

```bash
go test ./internal/helpers/... -v
```

### Manual testing

A Node.js client script is included for quick manual verification:

```bash
npm i
node client.js
```

This connects to the server, sends `ip.wav` over the WebSocket, and writes the converted output to `op.flac` in the project directory. Three sample audio files are included: `input.wav`, `ip.wav`, and `audio.wav`.

---

## Project Structure

```
Audio-Conversion-Servic/
├── cmd/                        # Application entrypoint
├── internal/
│   └── helpers/
│       ├── stream.go           # ProcessAudioStream — core conversion logic
│       └── stream_test.go      # Unit tests with mocked WebSocket + FFmpeg
├── tests/
│   └── server_integration_test.go  # End-to-end integration test
├── bin/
│   └── app                     # Pre-built binary
├── client.js                   # Manual test client (Node.js)
├── input.wav                   # Sample audio file
├── ip.wav                      # Sample audio file
├── audio.wav                   # Sample audio file
├── docker-compose.yml
├── Dockerfile
└── README.md
```

---

## Tech Stack

| Component | Technology | Purpose |
|---|---|---|
| **Language** | Go | Backend server, WebSocket handling, process management |
| **Audio Conversion** | FFmpeg | WAV to FLAC encoding via stdin/stdout piping |
| **Transport** | WebSockets (gorilla/websocket) | Bidirectional real-time streaming |
| **Containerization** | Docker + Docker Compose | Reproducible build and deployment |
| **Testing** | Go `testing` package | Integration and unit tests with mocks |

---

## FAQ's

**"Why WebSockets instead of a REST upload endpoint?"**
> A REST endpoint would require the client to upload the entire file, wait for conversion, then download the result — three separate steps. WebSockets let us stream the file in and the result out over a single persistent connection, reducing latency and simplifying the client.

**"Why FFmpeg instead of a Go audio library?"**
> FFmpeg is the industry standard for audio/video processing with decades of codec support and optimization. Piping through FFmpeg via stdin/stdout is simpler and more reliable than reimplementing codec logic in Go, and it supports virtually every audio format if you want to extend beyond WAV/FLAC.

**"Why is the FFmpeg command mocked in unit tests?"**
> The unit tests for `ProcessAudioStream` focus on the Go logic — WebSocket message handling, stream piping, and error propagation. Mocking FFmpeg isolates these concerns from the actual conversion binary, making tests fast, deterministic, and runnable without FFmpeg installed.

**"Can this handle large files?"**
> The streaming architecture means the entire file doesn't need to fit in memory — data is piped through FFmpeg as it arrives. That said, WebSocket message size limits and connection timeouts may need tuning for very large files in production.

---

## Disclaimer

This is a portfolio project demonstrating real-time audio processing with Go and WebSockets. It is not intended for production use without additional hardening — particularly around input validation, file size limits, connection timeouts, and authentication.

---

## License

MIT — see [LICENSE](LICENSE) for details.

---

<p align="center">
  <sub>Built by <a href="https://github.com/prranavv">prranavv</a></sub>
</p>
