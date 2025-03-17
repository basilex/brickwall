package auth

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
)

func Command(ctx context.Context) *cli.Command {
	command := &cli.Command{
		Name:     "auth",
		Category: "services",
		Usage:    "Run the auth service",
		Action: func(ctx context.Context, cli *cli.Command) error {
			return bootstrap(ctx, cli)
		},
	}
	return command
}

func bootstrap(_ context.Context, _ *cli.Command) error {
	log.Println("DEBUG: Auth command executed")

	return nil
}
