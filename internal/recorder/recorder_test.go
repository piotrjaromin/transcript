package recorder

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecorder(t *testing.T) {
	// Setup test samples
	testSamples := []float32{0.1, 0.2, 0.3, 0.4, 0.5}

	t.Run("basic recording flow", func(t *testing.T) {
		// Enable test mode
		EnableTestMode(testSamples)
		defer DisableTestMode()

		// Create temp file
		tmpFile, err := os.CreateTemp("", "test-recording-*.wav")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		rec := NewRecorder(tmpFile.Name())

		// Start recording
		err = rec.StartRecording()
		require.NoError(t, err)

		// Stop recording
		filePath, err := rec.StopRecording()
		require.NoError(t, err)
		assert.Equal(t, tmpFile.Name(), filePath)

		// Verify file was created
		_, err = os.Stat(filePath)
		assert.NoError(t, err)
	})

	t.Run("double start should fail", func(t *testing.T) {
		// Enable test mode
		EnableTestMode(testSamples)
		defer DisableTestMode()

		rec := NewRecorder("")
		err := rec.StartRecording()
		require.NoError(t, err)

		err = rec.StartRecording()
		assert.Error(t, err)
	})

	t.Run("stop without start should fail", func(t *testing.T) {
		rec := NewRecorder("")
		_, err := rec.StopRecording()
		assert.Error(t, err)
	})

	t.Run("get audio data", func(t *testing.T) {
		// Enable test mode
		EnableTestMode(testSamples)
		defer DisableTestMode()

		rec := NewRecorder("")
		err := rec.StartRecording()
		require.NoError(t, err)

		// Stop recording
		_, err = rec.StopRecording()
		require.NoError(t, err)

		data := rec.GetAudioData()
		assert.Equal(t, testSamples, data)
	})
}
