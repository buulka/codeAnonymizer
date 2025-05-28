package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"unicode/utf8"
)

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
