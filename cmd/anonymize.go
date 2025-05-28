package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

type Config struct {
	InputFile      string
	OutputFile     string
	Language       string
	NormalizedLang string
	KeepComments   bool
}

type FileValidator struct{}

func NewFileValidator() *FileValidator {
	return &FileValidator{}
}

func (fv *FileValidator) FileExists(path string) (exists bool, err error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (fv *FileValidator) IsGoFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".go"
}

func (fv *FileValidator) IsBinaryFile(path string) (isGo bool, err error) {
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

func (fv *FileValidator) IsFileEmpty(path string) (isEmpty bool, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.Size() == 0, nil
}

func (fv *FileValidator) ValidateGoFile(path string) (err error) {
	// Input file existence check
	exists, err := fv.FileExists(path)
	if err != nil {
		return fmt.Errorf("file check failed: %v", err)
	}
	if !exists {
		return fmt.Errorf("file %s does not exist", path)
	}

	// Input file type check
	if !fv.IsGoFile(path) {
		return fmt.Errorf("file %s is not a .go file", path)
	}

	// Empty file validation
	if empty, err := fv.IsFileEmpty(path); err != nil {
		return fmt.Errorf("failed to check file size: %w", err)
	} else if empty {
		return fmt.Errorf("file is empty: %s", path)
	}

	// Check if binary
	if isBinary, err := fv.IsBinaryFile(path); err == nil && isBinary {
		return fmt.Errorf("file appears to be binary: %s", path)
	}

	return nil
}

type LanguageValidator struct {
	allowedLanguages map[string]string
}

func NewLanguageValidator() *LanguageValidator {
	return &LanguageValidator{
		allowedLanguages: map[string]string{
			"go":     "go",
			"golang": "go",
		},
	}
}

func (lv *LanguageValidator) ValidateLanguage(language string) (string, error) {
	key := strings.ToLower(strings.TrimSpace(language))
	normalized, ok := lv.allowedLanguages[key]
	if !ok {
		return "", fmt.Errorf("unsupported language: %s", language)
	}
	return normalized, nil
}

type ConfigValidator struct {
	fileValidator     *FileValidator
	languageValidator *LanguageValidator
}

func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		fileValidator:     NewFileValidator(),
		languageValidator: NewLanguageValidator(),
	}
}

func (cv *ConfigValidator) Validate(config *Config) error {
	NormalizedLang, err := cv.languageValidator.ValidateLanguage(config.Language)
	if err != nil {
		log.Fatalf("Language error: %v", err)
	}
	config.NormalizedLang = NormalizedLang

	err = cv.fileValidator.ValidateGoFile(config.InputFile)
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

		validator := NewConfigValidator()
		if err := validator.Validate(&cfg); err != nil {
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

	if err := AnonymizeCmd.MarkFlagRequired("input"); err != nil {
		log.Printf("failed to mark input flag as required: %v", err)
	}

}
