package whisper

import (
	"fmt"
	"os"
	"strings"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

// Client interface for whisper transcription
type Client interface {
	Transcribe(samples []float32) (string, error)
	Close()
}

// WhisperClient implements the Client interface
type WhisperClient struct {
	model      whisper.Model
	context    whisper.Context
	language   string
	numThreads int
}

// NewClient creates a new whisper client
func NewClient(modelPath, language string, numThreads int) (Client, error) {
	if modelPath == "" {
		return nil, fmt.Errorf("model path is required")
	}

	// Check if model file exists
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model file not found: %s", modelPath)
	}

	// Load the model
	model, err := whisper.New(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}

	// Create context
	context, err := model.NewContext()
	if err != nil {
		model.Close()
		return nil, fmt.Errorf("failed to create context: %w", err)
	}

	// Set language if specified
	if language != "" && language != "auto" {
		if !model.IsMultilingual() {
			model.Close()
			return nil, fmt.Errorf("model is not multilingual but language '%s' was specified", language)
		}
		if err := context.SetLanguage(language); err != nil {
			model.Close()
			return nil, fmt.Errorf("unsupported language '%s' for this model: %v", language, err)
		}
	}

	// Set number of threads to use
	context.SetThreads(uint(numThreads))

	return &WhisperClient{
		model:      model,
		context:    context,
		language:   language,
		numThreads: numThreads,
	}, nil
}

// Close releases resources
func (c *WhisperClient) Close() {
	if c.model != nil {
		c.model.Close()
	}
}

// Transcribe transcribes audio data
func (c *WhisperClient) Transcribe(samples []float32) (string, error) {
	// Process the audio data
	err := c.context.Process(samples, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to process audio: %w", err)
	}

	// Build the transcript from all segments
	var segments []string
	for {
		segment, err := c.context.NextSegment()
		if err != nil {
			break // End of segments
		}
		segments = append(segments, segment.Text)
	}

	return strings.Join(segments, " "), nil
}
