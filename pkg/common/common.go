package common

import (
	"fmt"
	"github.com/urfave/cli/v3"
)

var Flags = []cli.Flag{
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

func ValidateCommonFlags(c *cli.Command) error {
	if c == nil {
		return fmt.Errorf("command cannot be nil")
	}
	if err := ValidateLogFormat(c.String("log-format")); err != nil {
		return fmt.Errorf("invalid log format: %s", err)
	}

	return nil
}
func ValidateLogFormat(logFormat string) error {
	switch logFormat {
	case "text", "json":
		return nil
	default:
		return fmt.Errorf("invalid log format: %s", logFormat)
	}
}
