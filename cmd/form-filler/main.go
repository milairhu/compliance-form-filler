package main

import (
	"compliance-form-filler/pkg/app"
	"compliance-form-filler/pkg/logger"
	"context"
	"os"
)

func main() {
	app := app.InitApp()

	if err := app.Run(context.Background(), os.Args); err != nil {
		logger.DefaultLogger.Fatal().Err(err).Msg("Failed to run")
	}
}
