package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Default model paths to check
var defaultModelPaths = []string{
	"./models/ggml-base.bin",
	"./models/ggml-small.bin",
	"./models/ggml-medium.bin",
	"./models/ggml-large.bin",
}

// For thread-safe model path resolution
var (
	modelOnce     sync.Once
	foundModelPath string
)

// getModelInfo returns a string describing the model being used
func getModelInfo() string {
	if modelPath != "" {
		return modelPath
	}
	
	// Find model path using thread-safe method
	path, _ := getModelPath()
	if path != "" {
		return fmt.Sprintf("default model found at %s", path)
	}
	
	return "no model found, please specify with --model flag"
}

// validateModelPath checks if a model path exists
func validateModelPath(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("specified model not found: %s", path)
	}
	return path, nil
}

// getModelPath returns the path to the model, or an error if no model is found
func getModelPath() (string, error) {
	if modelPath != "" {
		return validateModelPath(modelPath)
	}
	
	// Use sync.Once to ensure we only search for default models once
	modelOnce.Do(func() {
		for _, path := range defaultModelPaths {
			absPath, _ := filepath.Abs(path)
			if _, err := os.Stat(absPath); err == nil {
				foundModelPath = absPath
				return
			}
		}
	})
	
	if foundModelPath != "" {
		return foundModelPath, nil
	}
	
	return "", fmt.Errorf("no model found, please specify with --model flag")
}
