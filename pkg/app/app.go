package app

import (
	"compliance-form-filler/internal/answer"
	"compliance-form-filler/pkg/common"
	"github.com/urfave/cli/v3"
)

func InitApp() *cli.Command {
	app := &cli.Command{
		Name:  "compliance-form-filler",
		Usage: "EVERTRUST Compliance form Filler",
		Commands: []*cli.Command{
			answer.Command,
		},
		Flags: common.Flags,
	}
	return app
}
