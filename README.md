## WAV TO FLAC Audio Conversion Service(Backend)

## Instructions to Run the App

To run this application,first clone this repository onto your machine and ensure the port 8080 is freed.

Ensure the docker daemon is running and run `docker compose up --build`

## API Endpoints

**/stream** - This endpoint recieves a WAV file through a Websocket connection and converts it into a FLAC file and is sent to the client via the websocket connection.

## Tests Written

### Integration Test

In the tests directory,there is server_integration_test.go file that is an integration test which runs and passes all the tests successfully.

### Unit Tests

In the internal/helpers directory there is a stream_test.go file. It tests the ProcessAudioStream function that converts the WAV file to a FLAC file and sends it through the websocket connection.
