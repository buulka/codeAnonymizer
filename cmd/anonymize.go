package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

type Config struct {
	InputFile      string
	OutputFile     string
	Language       string
	NormalizedLang string
	KeepComments   bool
}

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

func ValidateLanguage(language string) (normalized string, err error) {
	var allowedLanguages = map[string]string{
		"go":     "go",
		"golang": "go"}

	key := strings.ToLower(strings.TrimSpace(language))
	normalized, ok := allowedLanguages[key]

	if !ok {
		return "", fmt.Errorf("unsupported language: %s", language)
	}
	return normalized, nil
}

func Validate(config *Config) error {
	NormalizedLang, err := ValidateLanguage(config.Language)
	if err != nil {
		log.Fatalf("Language error: %v", err)
	}
	config.NormalizedLang = NormalizedLang
	return nil
}

var AnonymizeCmd = &cobra.Command{
	Use:   "anonymize",
	Short: "Anonymize source code by replacing identifiers",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := Config{
			InputFile:    cmd.Flag("input").Value.String(),
			OutputFile:   cmd.Flag("output").Value.String(),
			Language:     cmd.Flag("language").Value.String(),
			KeepComments: cmd.Flag("keep-comments").Changed}

		if err := Validate(&cfg); err != nil {
			log.Fatalf("Configuration error: %v", err)
		}
		fmt.Println(cfg.NormalizedLang)

		// Input file existence check
		fileExists, err := DoesFileExist(cfg.InputFile)
		if err != nil {
			panic(err)
		}
		fmt.Println("\nInput file", cfg.InputFile, "exists:", fileExists, "\nLang:", cfg.NormalizedLang)

	},
}

func init() {
	rootCmd.AddCommand(AnonymizeCmd)

	AnonymizeCmd.Flags().StringP("input", "i", "", "Input source file (required)")
	AnonymizeCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")
	AnonymizeCmd.Flags().StringP("language", "l", "go", "Language to use")
	AnonymizeCmd.Flags().BoolP("keep-comments", "k", false, "Keep comments")

	err := AnonymizeCmd.MarkFlagRequired("input")
	if err != nil {
		return
	}

}
