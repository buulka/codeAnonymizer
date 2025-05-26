package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	InputFile    string
	OutputFile   string
	Language     string
	KeepComments bool
)

func DoesFileExist(path string) (found bool, err error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
	} else {
		found = true
	}
	return
}

var AnonymizeCmd = &cobra.Command{
	Use:   "anonymize",
	Short: "Anonymize source code by replacing identifiers",
	Run: func(cmd *cobra.Command, args []string) {

		if Language != "go" {
			log.Fatalf("Unsupported language: %s", Language)
		}

		// Input file existence check
		fileExists, err := DoesFileExist(InputFile)
		if err != nil {
			panic(err)
		}
		fmt.Println("\nInput file", InputFile, "exists:", fileExists)

	},
}

func init() {
	rootCmd.AddCommand(AnonymizeCmd)

	AnonymizeCmd.Flags().StringVarP(&InputFile, "input", "i", "", "Input source file (required)")
	AnonymizeCmd.Flags().StringVarP(&OutputFile, "output", "o", "", "Output file (default: stdout)")
	AnonymizeCmd.Flags().StringVarP(&Language, "language", "l", "go", "Language to use")
	AnonymizeCmd.Flags().BoolVarP(&KeepComments, "keep-comments", "k", false, "Keep comments")

	err := AnonymizeCmd.MarkFlagRequired("input")
	if err != nil {
		return
	}

}
