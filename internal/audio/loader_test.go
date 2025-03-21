package audio

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAudioFile(t *testing.T) {
	t.Run("valid WAV file", func(t *testing.T) {
		tempDir := os.TempDir()
		testFile := filepath.Join(tempDir, "test.wav")
		createValidTestWAV(t, testFile)
		defer os.Remove(testFile)

		samples, err := LoadAudioFile(testFile)
		require.NoError(t, err, "Valid WAV file should load without error")
		assert.NotEmpty(t, samples, "Should get audio samples from valid file")
	})

	t.Run("invalid file format", func(t *testing.T) {
		// Create temp text file
		tempDir := os.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		os.WriteFile(testFile, []byte("invalid data"), 0644)
		defer os.Remove(testFile)

		_, err := LoadAudioFile(testFile)
		assert.ErrorContains(t, err, "unsupported audio format")
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := LoadAudioFile("nonexistent.wav")
		assert.ErrorContains(t, err, "failed to open audio file")
	})
}

func createValidTestWAV(t *testing.T, path string) {
	file, err := os.Create(path)
	require.NoError(t, err)
	defer file.Close()

	// Create valid PCM data (1 second of silence at 16kHz)
	numSamples := 16000
	data := make([]int16, numSamples)
	header := make([]byte, 44)

	// RIFF header
	copy(header[0:4], "RIFF")
	binary.LittleEndian.PutUint32(header[4:8], uint32(36+len(data)*2)) // Chunk size
	copy(header[8:12], "WAVE")
	
	// fmt chunk
	copy(header[12:16], "fmt ")
	binary.LittleEndian.PutUint32(header[16:20], 16)    // Subchunk size
	binary.LittleEndian.PutUint16(header[20:22], 1)     // PCM format
	binary.LittleEndian.PutUint16(header[22:24], 1)     // Mono
	binary.LittleEndian.PutUint32(header[24:28], 16000) // Sample rate
	binary.LittleEndian.PutUint32(header[28:32], 32000) // Byte rate
	binary.LittleEndian.PutUint16(header[32:34], 2)     // Block align
	binary.LittleEndian.PutUint16(header[34:36], 16)    // Bits per sample
	
	// data chunk
	copy(header[36:40], "data")
	binary.LittleEndian.PutUint32(header[40:44], uint32(len(data)*2)) // Data size

	// Write header and data
	file.Write(header)
	binary.Write(file, binary.LittleEndian, data)
}
