package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	InputFile  string
	OutputFile string
	Language   string
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

func ParseFlags() Config {
	var config Config

	flag.StringVar(&config.InputFile, "input", "input.go", "an input file path")
	flag.StringVar(&config.OutputFile, "output", "output.go", "an output file path")
	flag.StringVar(&config.Language, "lang", "go", "file language")

	flag.Parse()

	fmt.Println("Args:")
	fmt.Println("\tinput file:", config.InputFile, "\n\toutput file:", config.OutputFile, "\n\tlang:", config.Language)

	return config
}

func main() {
	// CLI initialisation
	config := ParseFlags()

	// Input file existence check
	fileExists, err := DoesFileExist(config.InputFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nInput file", config.InputFile, "exists:", fileExists)
}
