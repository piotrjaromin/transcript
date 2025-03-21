# Transcript

Transcribe speech to text either by using CLI tool or starting HTTP server.

Project is using golang bindings for https://github.com/ggerganov/whisper.cpp.

## Features

- Convert audio files to text transcripts
- Support for multiple languages (including English and Polish)
- Three operation modes:
  1. HTTP server with API endpoint
  2. CLI tool for processing audio files
  3. CLI tool with live recording capability

## Installation

### Using Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/yourusername/transcript.git
cd transcript

# Build and run with Docker Compose
docker-compose up
```

### Manual Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/transcript.git
cd transcript

# Download a model (if not already present)
mkdir -p models
wget -O models/ggml-medium.en.bin https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-medium.en.bin

# Build the application
go build -o transcript

# Run the application
./transcript --help
```

## Usage

### HTTP Server Mode

Start the server:
```bash
# Basic startup with required model
./transcript server --port 8080 --model models/ggml-medium.en.bin

# With additional options
./transcript server \
  --port 8080 \
  --model models/ggml-medium.en.bin \
  --language en \
  --threads 4
```

API Endpoints:
- `POST /transcribe` - Upload audio file for transcription
  - Form parameters:
    - `audio` - Audio file in WAV format (max 10MB)

Example using curl:
```bash
curl -X POST -F "audio=@input.wav" http://localhost:8080/transcribe
```
Successful response will include the transcript and detected language.

### CLI File Transcription Mode

```bash
./transcript file --model models/ggml-medium.en.bin --input path/to/audio.wav
```

### CLI Recording Mode

```bash
./transcript record --model models/ggml-medium.en.bin
```

## Docker Configuration

You can specify a different model using the `WHISPER_MODEL` environment variable:

```bash
# Use a smaller model
WHISPER_MODEL=ggml-small.en.bin docker-compose up

# Use a larger model
WHISPER_MODEL=ggml-large-v3.bin docker-compose up
```

Available models:
- `ggml-tiny.en.bin` (39MB) - Fastest, least accurate
- `ggml-small.en.bin` (139MB) - Good balance for most uses
- `ggml-medium.en.bin` (1.5GB) - Default, good accuracy
- `ggml-large-v3.bin` (3.1GB) - Most accurate, slowest

## Development

### Prerequisites

- Go 1.21 or higher
- GCC compiler (for CGO)
- PortAudio (for recording functionality)

### Building from Source

```bash
go build -o transcript
```

### Running Tests

```bash
go test ./...
```

### CI/CD Pipeline

This project uses GitHub Actions for continuous integration and deployment:

- Automatically runs tests on pull requests
- Builds and pushes Docker images on tag pushes
- Publishes images to Docker Hub with appropriate tags

To trigger a new release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

Required GitHub secrets for deployment:
- `DOCKERHUB_USERNAME`: Your Docker Hub username
- `DOCKERHUB_TOKEN`: Docker Hub access token with write permissions

## License

[MIT License](LICENSE)
