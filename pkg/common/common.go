package common

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var Flags = []cli.Flag{
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
		Required: true,
		Value:    "",
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
	&cli.BoolFlag{
		Name:    "verbose",
		Sources: cli.EnvVars("VERBOSE"),
		Usage:   "Enable verbose logging",
		Local:   false,
		Value:   false,
	},
	&cli.StringFlag{
		Name:    "log-format",
		Sources: cli.EnvVars("LOG_FORMAT"),
		Usage:   "Log format (text or json)",
		Local:   false,
		Value:   "json",
	},
}

func ValidateFlags(c *cli.Command) error {
	if c == nil {
		return fmt.Errorf("command cannot be nil")
	}

	if c.String("source-file") == "" {
		return fmt.Errorf("source-file is required")
	}
	if c.String("output-file") == "" {
		return fmt.Errorf("output-file is required")
	}
	if !isValidFilePath(c.String("source-file")) {
		return fmt.Errorf("invalid source-file path: %s", c.String("source-file"))
	}
	if !isValidFilePath(c.String("output-file")) {
		return fmt.Errorf("invalid output-file path: %s", c.String("output-file"))
	}
	if err := ValidateLogFormat(c.String("log-format")); err != nil {
		return fmt.Errorf("invalid log format: %s", err)
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

func ValidateLogFormat(logFormat string) error {
	switch logFormat {
	case "text", "json":
		return nil
	default:
		return fmt.Errorf("invalid log format: %s", logFormat)
	}
}
