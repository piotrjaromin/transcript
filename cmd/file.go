package cmd

import (
	"fmt"

	"github.com/piotrjaromin/transcript/internal/transcriber"
	"github.com/spf13/cobra"
)

var (
	filePath string
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Transcribe an audio file",
	Long:  `Transcribe the specified audio file to text.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if filePath == "" {
			return fmt.Errorf("file path is required")
		}
		
		fmt.Printf("Transcribing file: %s\n", filePath)
		fmt.Printf("Using model: %s\n", getModelInfo())
		
		// Get the model path
		modelPath, err := getModelPath()
		if err != nil {
			return err
		}
		
		// Create a transcriber
		transcriber, err := transcriber.NewFileTranscriber(modelPath, language, numThreads)
		if err != nil {
			return fmt.Errorf("failed to create transcriber: %w", err)
		}
		defer transcriber.Close()
		
		// Transcribe the file
		transcript, err := transcriber.Transcribe(filePath)
		if err != nil {
			return fmt.Errorf("transcription failed: %w", err)
		}
		
		// Print the transcript
		fmt.Println("\nTranscript:")
		fmt.Println("----------")
		fmt.Println(transcript)
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
	
	fileCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the audio file to transcribe (required)")
	fileCmd.MarkFlagRequired("file")
}
