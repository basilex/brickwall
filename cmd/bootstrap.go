package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"brickwall/cmd/api"
	"brickwall/cmd/auth"
	"brickwall/internal/common"
)

func Bootstrap(ctx context.Context) error {
	md := ctx.Value(common.KeyMetadata).(*common.Metadata)

	version := fmt.Sprintf("%s-%s-%s", md.Version, md.Staging, md.Githash)

	command := &cli.Command{
		Name:      "Brickwall SaaS platform manager",
		Copyright: "Copyright (C) 2025 by Brickwall Inc. All Rights Reserved.",
		Version:   fmt.Sprintf("%s, %s", version, md.Gobuild),
		Usage:     "bsp <service> [flags]",

		Commands: []*cli.Command{
			api.Command(ctx),
			auth.Command(ctx),
		},
	}
	return command.Run(ctx, os.Args)
}
