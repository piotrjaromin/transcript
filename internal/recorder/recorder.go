package recorder

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/go-audio/audio"
	wavgo "github.com/go-audio/wav"
	audioloader "github.com/piotrjaromin/transcript/internal/audio"
)

// For testing purposes
var (
	// Allow tests to override the command execution
	execCommand = exec.Command
	// Flag to indicate if we're in test mode
	testMode bool
	// Test samples to return when in test mode
	testSamples []float32
)

// EnableTestMode turns on test mode for the recorder
func EnableTestMode(samples []float32) {
	testMode = true
	testSamples = samples
}

// DisableTestMode turns off test mode
func DisableTestMode() {
	testMode = false
	testSamples = nil
}

// Recorder handles audio recording functionality
type Recorder struct {
	outputFile string
	samples    []float32
	recording  bool
	sampleRate int
	cmd        *exec.Cmd
	tempFile   string
}

// NewRecorder creates a new audio recorder
func NewRecorder(outputFile string) *Recorder {
	return &Recorder{
		outputFile: outputFile,
		sampleRate: 16000, // Whisper expects 16kHz
		samples:    make([]float32, 0),
	}
}

// StartRecording starts recording audio using the system's default recording device
// This is a simplified implementation that uses sox for recording
func (r *Recorder) StartRecording() error {
	if r.recording {
		return fmt.Errorf("already recording")
	}

	// If we're in test mode, just set recording flag and return
	if testMode {
		r.recording = true
		return nil
	}

	// Create a temporary file for recording
	tempFile, err := os.CreateTemp("", "recording-*.wav")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempFile.Close()
	tempFileName := tempFile.Name()
	
	// Check if sox is installed
	_, err = exec.LookPath("sox")
	if err != nil {
		return fmt.Errorf("sox audio tool is required: %w", err)
	}
	
	// Start recording to the temporary file using sox in a separate process
	// This will continue until StopRecording is called
	cmd := execCommand("sox", "-d", "-r", fmt.Sprintf("%d", r.sampleRate), "-c", "1", tempFileName)
	
	// Start the recording process
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start recording: %w", err)
	}
	
	r.recording = true
	
	// Store the command process to stop it later
	r.cmd = cmd
	r.tempFile = tempFileName
	
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// StopRecording stops recording audio and returns the path to the recorded file
func (r *Recorder) StopRecording() (string, error) {
	if !r.recording {
		return "", fmt.Errorf("not currently recording")
	}

	r.recording = false

	// If we're in test mode, use the test samples
	if testMode {
		r.samples = testSamples
		if r.outputFile != "" {
			// Create a dummy file for testing
			if err := r.Save(r.outputFile); err != nil {
				return "", err
			}
		}
		return r.outputFile, nil
	}

	// Stop the recording process
	if r.cmd != nil && r.cmd.Process != nil {
		// Send SIGTERM to the sox process
		r.cmd.Process.Signal(syscall.SIGTERM)
		
		// Wait for the process to exit
		r.cmd.Wait()
	}

	// Check if the temporary file exists
	if _, err := os.Stat(r.tempFile); os.IsNotExist(err) {
		return "", fmt.Errorf("recording file not found, recording may have failed")
	}

	// Load the recorded audio file
	samples, err := audioloader.LoadAudioFile(r.tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to load recorded audio: %w", err)
	}
	
	// Store the samples
	r.samples = samples

	// Save to the output file if specified
	if r.outputFile != "" {
		err = copyFile(r.tempFile, r.outputFile)
		if err != nil {
			return "", fmt.Errorf("failed to save recording to output file: %w", err)
		}
		// Clean up the temporary file after successful copy
		os.Remove(r.tempFile)
	} else {
		// Clean up the temporary file if we're not keeping it
		os.Remove(r.tempFile)
	}

	return r.outputFile, nil
}

// GetAudioData returns the recorded audio data
func (r *Recorder) GetAudioData() []float32 {
	return r.samples
}

// Save saves the recorded audio to a file
func (r *Recorder) Save(path string) error {
	if len(r.samples) == 0 {
		return fmt.Errorf("no audio data recorded")
	}

	// Create the output file
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Create a new encoder
	enc := wavgo.NewEncoder(f, r.sampleRate, 16, 1, 1) // 1 is for PCM format
	
	// Convert float32 samples to int
	intSamples := make([]int, len(r.samples))
	for i, sample := range r.samples {
		// Convert normalized float32 [-1.0,1.0] to int16 range
		intSamples[i] = int(sample * 32767)
	}
	
	// Create audio.IntBuffer
	buf := &audio.IntBuffer{
		Data:           intSamples,
		Format:         &audio.Format{SampleRate: r.sampleRate, NumChannels: 1},
		SourceBitDepth: 16,
	}
	
	// Write the buffer to the file
	if err := enc.Write(buf); err != nil {
		return fmt.Errorf("failed to write audio data: %w", err)
	}
	
	// Close the encoder
	if err := enc.Close(); err != nil {
		return fmt.Errorf("failed to close encoder: %w", err)
	}

	return nil
}
