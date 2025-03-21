package transcriber

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileTranscriber_TranscribeFromSamples(t *testing.T) {
	mockClient := &mockWhisperClient{
		transcribeFunc: func(samples []float32) (string, error) {
			return "test transcription", nil
		},
	}

	transcriber := &FileTranscriber{
		modelPath: "test-model",
		language:  "en",
		client:    mockClient,
		newClient: func(modelPath, language string, threads int) (whisperClient, error) {
			return mockClient, nil
		},
	}

	// Test transcription
	samples := []float32{0.1, 0.2, 0.3}
	transcript, err := transcriber.TranscribeFromSamples(samples)
	require.NoError(t, err)
	assert.Equal(t, "test transcription", transcript)
}

type mockWhisperClient struct {
	transcribeFunc func([]float32) (string, error)
	closeFunc      func()
}

func (m *mockWhisperClient) Transcribe(samples []float32) (string, error) {
	return m.transcribeFunc(samples)
}

func (m *mockWhisperClient) Close() {
	if m.closeFunc != nil {
		m.closeFunc()
	}
}
