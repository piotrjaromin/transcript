# Stage 1: Builder
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git wget

WORKDIR /app

# Create models directory
RUN mkdir -p models

# Download default model
RUN wget -O models/ggml-medium.en.bin https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-medium.en.bin

# Copy go mod and sum files
COPY go.mod go.sum* ./

# Download Go dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with version info
ARG VERSION=dev
ARG COMMIT_SHA=unknown
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-X 'github.com/piotrjaromin/transcript/internal/version.Version=${VERSION}' -X 'github.com/piotrjaromin/transcript/internal/version.CommitSHA=${COMMIT_SHA}'" -o transcript

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates portaudio alsa-lib wget

WORKDIR /app

# Create models directory
RUN mkdir -p models

# Copy built binary
COPY --from=builder /app/transcript /app/transcript
# Copy default model
COPY --from=builder /app/models /app/models

# Entrypoint script for model handling
RUN echo $'#!/bin/sh\n\
if [ -n "$WHISPER_MODEL" ]; then\n\
    if [ ! -f "/app/models/$WHISPER_MODEL" ]; then\n\
        echo "Downloading custom model: $WHISPER_MODEL"\n\
        wget -O "/app/models/$WHISPER_MODEL" "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/$WHISPER_MODEL"\n\
    fi\n\
    MODEL_PATH="/app/models/$WHISPER_MODEL"\n\
else\n\
    MODEL_PATH="/app/models/ggml-medium.en.bin"\n\
fi\n\
\n\
exec /app/transcript --model "$MODEL_PATH" "$@"' > entrypoint.sh && \
    chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
