package answer

import (
	"compliance-form-filler/pkg/common"

	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
)

var Command = &cli.Command{
	Name:  "answer",
	Usage: "Answer the questions in the source file and save the questions/answers to the destination file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "source-file",
			Usage:    ".txt file containing the questions",
			Sources:  cli.EnvVars("SOURCE_FILE"),
			Required: true,
			Value:    "",
		},
		&cli.StringFlag{
			Name:     "output-file",
			Usage:    ".csv file to save the answers to questions",
			Sources:  cli.EnvVars("OUTPUT_FILE"),
			Required: false,
			Value:    "results.csv",
		},
		&cli.StringFlag{
			Name:    "qdrant-url",
			Usage:   "Qdrant URL for vector database",
			Sources: cli.EnvVars("QDRANT_URL"),
			Value:   "http://localhost:6333",
		},
		&cli.StringFlag{
			Name:     "llm-url",
			Usage:    "URL for the LLM service",
			Sources:  cli.EnvVars("LLM_URL"),
			Required: true,
			Value:    "",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return validateAndExecute(cmd)
	},
}

func ValidateFlags(cmd *cli.Command) error {
	if cmd.String("source-file") == "" {
		return fmt.Errorf("source-file is required")
	}
	if cmd.String("output-file") == "" {
		return fmt.Errorf("output-file is required")
	}
	if cmd.String("qdrant-url") == "" {
		return fmt.Errorf("qdrant-url is required")
	}
	if cmd.String("llm-url") == "" {
		return fmt.Errorf("llm-url is required")
	}
	if !isValidFilePath(cmd.String("source-file")) {
		return fmt.Errorf("invalid source-file path: %s", cmd.String("source-file"))
	}
	if !checkFileExtension(cmd.String("source-file"), ".txt") {
		return fmt.Errorf("source-file must be a .txt file: %s", cmd.String("source-file"))
	}
	if !checkFileExtension(cmd.String("output-file"), ".csv") {
		return fmt.Errorf("output-file must be a .csv file: %s", cmd.String("output-file"))
	}
	return nil
}

func isValidFilePath(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

func checkFileExtension(path string, ext string) bool {
	if len(path) < len(ext) {
		return false
	}
	return path[len(path)-len(ext):] == ext
}

func validateAndExecute(cmd *cli.Command) error {
	// Validate global flags
	if err := common.ValidateCommonFlags(cmd); err != nil {
		return err
	}

	// Validate specific flags for this command
	if err := ValidateFlags(cmd); err != nil {
		return err
	}

	return Answer(cmd)
}
