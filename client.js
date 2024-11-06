const fs = require('fs');
const WebSocket = require('ws');
const path = require('path');

const wavFilePath = path.join(__dirname, 'ip.wav');
const flacFilePath = path.join(__dirname, 'op.flac');
const flacWriteStream = fs.createWriteStream(flacFilePath);

const ws = new WebSocket('http://localhost:8080/stream');

ws.on('open', () => {
    console.log('WebSocket connection opened.');
    sendWAVFile();
});
// Function to read and send the WAV file
function sendWAVFile() {
    const wavStream = fs.createReadStream(wavFilePath);
    let isHeaderSent = false;
    let headerBuffer = Buffer.alloc(44); // WAV header is typically 44 bytes
    wavStream.on('data', (chunk) => {
        if (!isHeaderSent) {
            if (chunk.length >= 44) {
                headerBuffer = chunk.slice(0, 44);
                ws.send(headerBuffer); // Send the header first
                console.log('WAV header sent.');
                isHeaderSent = true;
                if (chunk.length > 44) {
                    const audioData = chunk.slice(44);
                    ws.send(audioData);
                }
            } else {
                console.error('Chunk is smaller than expected header size.');
            }
        } else {
            ws.send(chunk);
        }
    });
    wavStream.on('end', () => {
        console.log('WAV file streaming completed.');

    });
    wavStream.on('error', (err) => {
        console.error('Error reading WAV file:', err);
        ws.close();
    });
}

ws.on('message', (data) => {
    if (Buffer.isBuffer(data)) {
        flacWriteStream.write(data);   
    }
});
