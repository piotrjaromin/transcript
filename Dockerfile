FROM golang:1.23-bookworm AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    gcc g++ git wget make libstdc++-12-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Create models directory
RUN mkdir -p models

COPY go.mod go.sum* ./

RUN go mod download

# Copy source code
COPY . .

# Build the application with version info
ARG VERSION=dev
ARG COMMIT_SHA=unknown

RUN make build

# Stage 2: Runtime
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates portaudio19-dev alsa-utils wget ffmpeg && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

ARG DEFAULT_MODEL=ggml-large-v3-turbo.bin

RUN mkdir -p models
RUN wget -O models/${DEFAULT_MODEL} "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/${DEFAULT_MODEL}?download=true"

# Copy built binary
COPY --from=builder /app/bin/transcript /app/transcript

# Set environment variable so that it persists
ENV WHISPER_MODEL=${DEFAULT_MODEL}

EXPOSE 8080

COPY ./scripts/entrypoint.sh /app/

RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]
