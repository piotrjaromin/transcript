package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// MoveFile safely moves files across filesystems
func MoveFile(src, dst string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Open source file
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	// Create destination file
	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	// Copy contents
	if _, err := io.Copy(destination, source); err != nil {
		return fmt.Errorf("failed to copy contents: %w", err)
	}

	// Remove source after successful copy
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("failed to remove source file: %w", err)
	}

	return nil
}
