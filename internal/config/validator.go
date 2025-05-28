package config

import (
	"fmt"
)

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

func (cv *ConfigValidator) Validate(cfg *Config) error {
	normalizedLang, err := cv.languageValidator.ValidateLanguage(cfg.Language)
	if err != nil {
		return fmt.Errorf("language validation failed: %w", err)
	}
	cfg.NormalizedLang = normalizedLang

	if err := cv.fileValidator.ValidateGoFile(cfg.InputFile); err != nil {
		return fmt.Errorf("file validation failed: %w", err)
	}

	return nil
}
