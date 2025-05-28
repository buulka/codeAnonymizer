package config

import (
	"fmt"
	"strings"
)

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
