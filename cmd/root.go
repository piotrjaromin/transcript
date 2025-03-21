package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	modelPath  string
	language   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "transcript",
	Short: "A tool for transcribing audio to text",
	Long: `Transcript is a CLI application that can transcribe audio files to text.
It can run in three modes:
1. HTTP server mode - serves an API endpoint for transcription
2. File mode - transcribes a specified audio file
3. Recorder mode - records audio from microphone and transcribes it`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return err
	}
	return nil
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&modelPath, "model", "", "Path to the whisper model file (if not provided, will use embedded model)")
	rootCmd.PersistentFlags().StringVar(&language, "language", "auto", "Language of the audio (optional, auto-detected if not provided)")
}
