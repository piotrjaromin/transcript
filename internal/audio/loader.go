package audio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unsafe"

	ffmpeg "github.com/u2takey/ffmpeg-go"
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
	
	return convertAudioWithFFmpeg(file)
}

// LoadAudioFromReader loads audio from an io.Reader
func LoadAudioFromReader(reader io.Reader) ([]float32, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}
	return convertAudioWithFFmpeg(bytes.NewReader(data))
}

// convertAudioWithFFmpeg converts audio from any format to float32 samples
// using FFmpeg for maximum compatibility with different audio formats
func convertAudioWithFFmpeg(input io.Reader) ([]float32, error) {
	var buf bytes.Buffer
	
	err := ffmpeg.Input("pipe:0").
		Output("pipe:1", ffmpeg.KwArgs{
			"f":           "f32le",
			"ar":          SampleRate,
			"ac":          1,
			"loglevel":    "error",
			"hide_banner": "",
		}).
		WithInput(input).
		WithOutput(&buf).
		Run()
	
	if err != nil {
		return nil, fmt.Errorf("ffmpeg conversion failed: %w", err)
	}

	// Convert byte buffer to float32 samples
	raw := buf.Bytes()
	if len(raw)%4 != 0 {
		return nil, fmt.Errorf("invalid f32le byte length: %d", len(raw))
	}

	samples := make([]float32, len(raw)/4)
	for i := 0; i < len(samples); i++ {
		// Convert little-endian bytes to float32
		bytes := raw[i*4 : i*4+4]
		bits := uint32(bytes[0]) | uint32(bytes[1])<<8 | uint32(bytes[2])<<16 | uint32(bytes[3])<<24
		samples[i] = *(*float32)(unsafe.Pointer(&bits))
	}

	return samples, nil
}
