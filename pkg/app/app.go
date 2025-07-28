package app

import (
	"compliance-form-filler/pkg/common"
	"compliance-form-filler/pkg/logger"
	"context"
	"github.com/urfave/cli/v3"
)

func InitApp() *cli.Command {
	app := &cli.Command{
		Name:  "compliance-form-filler",
		Usage: "EVERTRUST Compliance form Filler",
		Flags: common.Flags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			logger.NewFromCliContext(cmd)
			return fillForm(cmd)
		},
	}

	return app
}

func fillForm(cmd *cli.Command) error {

}
