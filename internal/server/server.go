package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/piotrjaromin/transcript/internal/audio"
	"github.com/piotrjaromin/transcript/internal/whisper"
)

// Server represents the HTTP server for transcription
type Server struct {
	port          int
	modelPath     string
	language      string
	numThreads    int
	whisperClient whisper.Client
}

// NewServer creates a new transcription server
func NewServer(port int, modelPath, language string, threads int) *Server {
	return &Server{
		port:       port,
		modelPath:  modelPath,
		language:   language,
		numThreads: threads,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	if s.port < 1 || s.port > 65535 {
		return fmt.Errorf("invalid port number: %d", s.port)
	}

	// Check if model file exists
	if _, err := os.Stat(s.modelPath); os.IsNotExist(err) {
		return fmt.Errorf("model file not found: %s", s.modelPath)
	}

	// Initialize whisper client
	var err error
	s.whisperClient, err = whisper.NewClient(s.modelPath, s.language, s.numThreads)
	if err != nil {
		return fmt.Errorf("failed to initialize whisper client: %w", err)
	}
	defer s.whisperClient.Close()

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20  // 8 MB limit for uploaded files

	r.POST("/transcribe", s.handleTranscribe)

	return r.Run(":" + strconv.Itoa(s.port))
}

// handleTranscribe handles the transcription endpoint
func (s *Server) handleTranscribe(c *gin.Context) {
	// Get audio file from request
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing audio file"})
		return
	}

	// Check file size
	if file.Size > 10*1024*1024 { // 10MB limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 10MB)"})
		return
	}

	// Create a secure temporary file
	tempFile, err := ioutil.TempFile("", "audio-*.wav")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp file"})
		return
	}
	tempFile.Close()
	
	// Save uploaded file to temp location
	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		os.Remove(tempFile.Name())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save audio file"})
		return
	}
	defer os.Remove(tempFile.Name())

	// Load audio samples
	samples, err := audio.LoadAudioFile(tempFile.Name())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid audio file: %v", err),
		})
		return
	}

	// Transcribe audio
	transcript, err := s.whisperClient.Transcribe(samples)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Transcription failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transcript": transcript,
		"language":   s.language,
	})
}
