# Code Anonymization Utility

A tool to anonymize source code by replacing variable names, functions, and other identifiers with generic placeholders (e.g., var1, func2), while optionally removing or obfuscating comments


Build a binary file

- `go build -o anonymizer`

Run 

- ` ./anonymizer anonymize -i input.go -o output.go -l go`
