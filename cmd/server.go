package cmd

import (
	"fmt"
	"os"

	"github.com/piotrjaromin/transcript/internal/server"
	"github.com/spf13/cobra"
)

var (
	port int
)

var numThreads = 4

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run as HTTP server",
	Long:  `Start an HTTP server that provides an API endpoint for audio transcription.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if model file exists
		if _, err := os.Stat(modelPath); err != nil {
			return fmt.Errorf("invalid model path: %w", err)
		}

		fmt.Printf("Starting HTTP server on port %d\n", port)
		fmt.Printf("Using model: %s\n", modelPath)
		fmt.Printf("Default language: %s\n", language)

		srv := server.NewServer(port, modelPath, language, numThreads)
		return srv.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&port, "port", 8080, "Port to run the HTTP server on")
	serverCmd.Flags().StringVar(&modelPath, "model", "", "Path to the whisper model file (required)")
	serverCmd.MarkFlagRequired("model")
}
