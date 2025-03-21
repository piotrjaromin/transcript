package transcriber

import (
	"fmt"
	"sync"

	"github.com/piotrjaromin/transcript/internal/audio"
	"github.com/piotrjaromin/transcript/internal/whisper"
)

// whisperClient defines the interface for whisper clients
type whisperClient interface {
	Transcribe(samples []float32) (string, error)
	Close()
}

// FileTranscriber handles transcription of audio files
type FileTranscriber struct {
	mu        sync.Mutex
	modelPath string
	language  string
	client    whisperClient
	newClient func(modelPath, language string, threads int) (whisperClient, error)
}

// NewFileTranscriber creates a new file transcriber
func NewFileTranscriber(modelPath, language string, threads int) (*FileTranscriber, error) {
	clientFactory := func(modelPath, language string, threads int) (whisperClient, error) {
		return whisper.NewClient(modelPath, language, threads)
	}
	
	// Create a client that will be reused for all transcriptions
	client, err := clientFactory(modelPath, language, threads)
	if err != nil {
		return nil, fmt.Errorf("failed to create whisper client: %w", err)
	}
	
	return &FileTranscriber{
		modelPath: modelPath,
		language:  language,
		client:    client,
		newClient: clientFactory,
	}, nil
}

// Close releases resources used by the transcriber
func (t *FileTranscriber) Close() {
	if t.client != nil {
		t.client.Close()
	}
}

// TranscribeFromSamples transcribes audio samples
func (t *FileTranscriber) TranscribeFromSamples(samples []float32) (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	// Use the existing client
	return t.client.Transcribe(samples)
}

// Transcribe transcribes the audio file at the given path
func (t *FileTranscriber) Transcribe(filePath string) (string, error) {
	// Load the audio file
	samples, err := audio.LoadAudioFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to load audio file: %w", err)
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	
	// Use the existing client
	transcript, err := t.client.Transcribe(samples)
	if err != nil {
		return "", fmt.Errorf("failed to transcribe audio: %w", err)
	}

	return transcript, nil
}
