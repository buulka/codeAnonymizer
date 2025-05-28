package cmd

import (
	config "anonimCode/internal/config"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var AnonymizeCmd = &cobra.Command{
	Use:   "anonymize",
	Short: "Anonymize source code by replacing identifiers",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.Config{
			InputFile:    cmd.Flag("input").Value.String(),
			OutputFile:   cmd.Flag("output").Value.String(),
			Language:     cmd.Flag("language").Value.String(),
			KeepComments: cmd.Flag("keep-comments").Changed}

		validator := config.NewConfigValidator()
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
