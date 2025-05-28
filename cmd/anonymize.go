package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type Config struct {
	InputFile      string
	OutputFile     string
	Language       string
	NormalizedLang string
	KeepComments   bool
}

func FileExists(path string) (exists bool, err error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func IsGoFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".go"
}

func IsBinaryFile(path string) (isGo bool, err error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	buf := make([]byte, 512)
	if _, err = file.Read(buf); err != nil {
		return false, err
	}

	return !utf8.Valid(buf), nil
}

func IsFileEmpty(path string) (isEmpty bool, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.Size() == 0, nil
}

func ValidateGoFile(path string) (err error) {
	// Input file existence check
	exists, err := FileExists(path)
	if err != nil {
		return fmt.Errorf("file check failed: %v", err)
	}
	if !exists {
		return fmt.Errorf("file %s does not exist", path)
	}

	// Input file type check
	if !IsGoFile(path) {
		return fmt.Errorf("file %s is not a .go file", path)
	}

	// Empty file validation
	if empty, err := IsFileEmpty(path); err != nil {
		return fmt.Errorf("failed to check file size: %w", err)
	} else if empty {
		return fmt.Errorf("file is empty: %s", path)
	}

	// Check if binary
	if isBinary, err := IsBinaryFile(path); err == nil && isBinary {
		return fmt.Errorf("file appears to be binary: %s", path)
	}

	return nil
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

	err = ValidateGoFile(config.InputFile)
	if err != nil {
		log.Fatalf("Input file type error: %v", err)
	}

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
