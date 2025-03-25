package api

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v3"

	"brickwall/cmd/api/service"
	"brickwall/internal/common"
	"brickwall/internal/provider"
)

func Command(ctx context.Context) *cli.Command {
	command := &cli.Command{
		Name:     "api",
		Category: "services",
		Usage:    "Run the api service",
		Action: func(ctx context.Context, cli *cli.Command) error {
			slog.Info("bsp: starting service bsp-api...")
			return bootstrap(ctx)
		},
	}
	return command
}

// @title       Brickwall API
// @version     0.1.0
// @description This is Brickwall RestAPI
// @host        localhost:8081
// @BasePath    /api/v1
func bootstrap(ctx context.Context) error {
	//
	// Logger provider - no dependencies
	//
	slog.SetDefault(
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
	)
	//
	// envProvider - no dependencies
	//
	envProvider := provider.NewEnvProvider(ctx)
	ctx = context.WithValue(ctx, common.KeyEnvProvider, envProvider)
	//
	// NATS provider - no dependencies
	//
	natsProvider := provider.NewNatsProvider(ctx)
	if _, err := natsProvider.Connect(); err != nil {
		return err
	}
	ctx = context.WithValue(ctx, common.KeyNatsProvider, natsProvider)
	defer natsProvider.Disconnect()
	//
	// Redis provider - no dependencies
	//
	redisProvider := provider.NewRedisProvider(ctx)
	if _, err := redisProvider.Connect(); err != nil {
		return err
	}
	ctx = context.WithValue(ctx, common.KeyRedisProvider, redisProvider)
	defer redisProvider.Disconnect()
	//
	// Pgx provider - no dependencies
	//
	pgxProvider := provider.NewPgxProvider(ctx)
	if _, err := pgxProvider.Connect(); err != nil {
		return err
	}
	ctx = context.WithValue(ctx, common.KeyPgxProvider, pgxProvider)
	defer pgxProvider.Disconnect()
	//
	// Jwt provider - depends on Redis
	//
	jwtProvider := provider.NewJwtProvider(ctx)
	ctx = context.WithValue(ctx, common.KeyJwtProvider, jwtProvider)
	//
	// twoFA provider - no dependencies
	//
	twoFAProvider := provider.New2FAProvider(ctx)
	ctx = context.WithValue(ctx, common.Key2FAProvider, twoFAProvider)
	//
	// Router provider - no dependencies
	//
	routerProvider := provider.NewRouterProvider(ctx).Init()
	ctx = context.WithValue(ctx, common.KeyRouterProvider, routerProvider)
	//
	// Validator provider - no dependencies
	//
	validator := validator.New()
	ctx = context.WithValue(ctx, common.KeyValidatorProvider, validator)
	//
	// SMTP provider - no dependencies
	//
	smtpProvider := provider.NewSmtpProvider(ctx)
	ctx = context.WithValue(ctx, common.KeySmtpProvider, smtpProvider)
	//
	// Service Manager - depends on pgx and sqlc storage.queries
	//
	serviceManager := service.NewServiceManager(ctx)
	ctx = context.WithValue(ctx, common.KeyServiceManager, serviceManager)
	//
	// Server provider - no dependencies
	// Here is HTTP(S) server starts in goroutine
	//
	RegisterRoutes(ctx, routerProvider)
	srv := provider.NewServerProvider(ctx).Startup(routerProvider)
	defer srv.Shutdown()
	//
	// Watcher provider - no dependencies
	// Here the app is waiting for the signals to interruption
	//
	provider.NewWatcherProvider().Catch()
	return nil
}
