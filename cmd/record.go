package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/piotrjaromin/transcript/internal/recorder"
	"github.com/piotrjaromin/transcript/internal/transcriber"
	"github.com/spf13/cobra"
)

var (
	outputFile string
)

// recordCmd represents the record command
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record and transcribe audio",
	Long:  `Record audio from the microphone, then transcribe it to text and printout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get model path
		modelPath, err := getModelPath()
		if err != nil {
			return err
		}

		fmt.Printf("Using model: %s\n", modelPath)

		// Create recorder
		rec := recorder.NewRecorder(outputFile)

		fmt.Println("Press ENTER to start recording...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		// Start recording
		err = rec.StartRecording()
		if err != nil {
			return fmt.Errorf("failed to start recording: %w", err)
		}

		fmt.Println("Recording... Press ENTER again to stop recording and start transcription.")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		// Stop recording
		_, err = rec.StopRecording()
		if err != nil {
			return fmt.Errorf("failed to stop recording: %w", err)
		}

		// Get audio data
		samples := rec.GetAudioData()
		if len(samples) == 0 {
			return fmt.Errorf("no audio data recorded")
		}

		// Create transcriber
		trans, err := transcriber.NewFileTranscriber(modelPath, language, numThreads)
		if err != nil {
			return fmt.Errorf("failed to create transcriber: %w", err)
		}
		defer trans.Close()

		// Transcribe audio
		fmt.Println("Transcribing audio...")
		transcript, err := trans.TranscribeFromSamples(samples)
		if err != nil {
			return fmt.Errorf("transcription failed: %w", err)
		}

		// Print transcript
		fmt.Println("\nTranscript:")
		fmt.Println("----------")
		fmt.Println(transcript)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
	recordCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Path to save the recorded audio (optional)")
}
