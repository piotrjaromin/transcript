package audio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-audio/wav"
)

// SampleRate is the sample rate expected by Whisper
const SampleRate = 16000

// LoadAudioFile loads an audio file and returns the samples as float32 values
func LoadAudioFile(filePath string) ([]float32, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".wav":
		return loadWAV(file)
	default:
		return nil, fmt.Errorf("unsupported audio format: %s", ext)
	}
}

// LoadAudioFromReader loads audio from an io.Reader
func LoadAudioFromReader(reader io.Reader) ([]float32, error) {
	// Convert to bytes first since wav.NewDecoder requires io.ReadSeeker
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}
	
	return loadWAVFromBytes(data)
}

// loadWAV loads a WAV file and returns the samples as float32 values
func loadWAV(file io.Reader) ([]float32, error) {
	// Convert to bytes first since wav.NewDecoder requires io.ReadSeeker
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}
	
	return loadWAVFromBytes(data)
}

// loadWAVFromBytes loads WAV data from a byte slice
func loadWAVFromBytes(data []byte) ([]float32, error) {
	// Create a bytes.Reader which implements io.ReadSeeker
	reader := bytes.NewReader(data)
	decoder := wav.NewDecoder(reader)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid WAV file")
	}

	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to read WAV data: %w", err)
	}

	// Convert int samples to float32 in range [-1, 1]
	samples := make([]float32, buffer.NumFrames())
	for i := 0; i < buffer.NumFrames(); i++ {
		samples[i] = float32(buffer.Data[i]) / 32768.0
	}

	// Check if audio is mono
	if buffer.Format.NumChannels > 1 {
		return nil, fmt.Errorf("only mono audio supported, got %d channels", 
			buffer.Format.NumChannels)
	}

	// Check if we need to resample
	if buffer.Format.SampleRate != SampleRate {
		fmt.Printf("Warning: audio sample rate is %dHz, resampling to %dHz\n", 
			buffer.Format.SampleRate, SampleRate)
		
		// Simple linear resampling (for proper implementation, use a resampling library)
		ratio := float64(SampleRate) / float64(buffer.Format.SampleRate)
		newLength := int(float64(len(samples)) * ratio)
		resampled := make([]float32, newLength)
		
		for i := 0; i < newLength; i++ {
			srcIdx := float64(i) / ratio
			idx1 := int(srcIdx)
			idx2 := idx1 + 1
			frac := float32(srcIdx - float64(idx1))
			
			if idx2 >= len(samples) {
				resampled[i] = samples[idx1]
			} else {
				resampled[i] = samples[idx1]*(1-frac) + samples[idx2]*frac
			}
		}
		
		samples = resampled
	}

	return samples, nil
}
