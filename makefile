CURRENT_DIR := $(shell pwd)

CGO_CFLAGS := -I${CURRENT_DIR}/cpplib
CGO_LDFLAGS := -L${CURRENT_DIR}/cpplib -lwhisper -lm -lstdc++ -fopenmp

export CGO_CFLAGS
export CGO_LDFLAGS

install:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod download

run:
	go run main.go --model ./models/ggml-large-v2-turbo.bin

run-file:
	go run main.go file --file audio.wav --model ./models/ggml-large-v2-turbo.bin --language pl 


run-record:
	go run main.go record --model ./models/ggml-large-v2-turbo.bin --language pl --output out.wav

run-server:
	go run main.go --model ./models/ggml-large-v2-turbo.bin server 

build:
	go build -o bin/transcript main.go

test: install
	go test ./...

