## WAV TO FLAC Audio Conversion Service(Backend)

## Instructions to Run the App

To run this application,first clone this repository onto your machine and ensure the port 8080 is freed.

Ensure the docker daemon is running and run `docker compose up --build`

If you dont want to run it through docker,you can run the binary present in the repository. Run the following command `bin/app` to run the app locally.

## API Endpoints

**/stream** - This endpoint recieves a WAV file through a Websocket connection and converts it into a FLAC file and is sent to the client via the websocket connection.

## Tests Written

### Integration Test

In the tests directory,there is server_integration_test.go file that is an integration test which runs and passes all the tests successfully.

### Unit Tests

In the internal/helpers directory there is a stream_test.go file. It tests the ProcessAudioStream function that converts the WAV file to a FLAC file and sends it through the websocket connection.

In this test,the Websocket connection,the execution command for `ffmpeg` are all mocked so all edge cases are taken care of.

### Manual Testing

To test the service manually I have written a client.js file that mocks a client sending a WAV file to the server through a Websocket connection.Make sure you have a Javascript runtime on your machine like node,deno or bun to run the JS file.

Run `npm i` to install the dependencies for the client.

In the case you have Node on your machine,make sure you are in the root directory and run `node client.js`. This fires up a Websocket client and connects to the server and sends a WAV file called `ip.wav` present in the root directory and in a matter of seconds we can see a flac file called `op.flac` present in the directory indicating the audio conversion was successfull.
