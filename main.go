package main

import (
	"flag"
	"fmt"
	"os"
)

func doesFileExist(path string) (found bool, err error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
	} else {
		found = true
	}
	return
}

func initCLI() (*string, *string, *string) {
	input := flag.String("input", "input.go", "an input file path")
	output := flag.String("output", "output.go", "an output file path")
	lang := flag.String("lang", "go", "file language")

	flag.Parse()

	fmt.Println("Args:")
	fmt.Println("\tinput file:", *input, "\n\toutput file:", *output, "\n\tlang:", *lang)

	return input, output, lang
}
func main() {
	// CLI initialisation
	input, _, _ := initCLI()

	// Input file existence check
	fileExists, err := doesFileExist(*input)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nInput file", *input, "exists:", fileExists)
}
