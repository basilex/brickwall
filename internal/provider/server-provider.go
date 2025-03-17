package provider

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/urfave/cli/v3"

	"brickwall/internal/common"
)

type IServerProvider interface {
	Startup(IRouterProvider) IServerProvider
	Shutdown() error
}

type ServerProvider struct {
	ctx    context.Context
	server *http.Server
}

func NewServerProvider(ctx context.Context) IServerProvider {
	return &ServerProvider{ctx: ctx}
}

func (rcv *ServerProvider) Startup(r IRouterProvider) IServerProvider {
	cli := rcv.ctx.Value(common.KeyCommand).(*cli.Command)

	go func() {
		slog.Info(
			"http/s server started", "bind", cli.String("server-address"),
		)
		if cli.Bool("tls-ssl-enabled") {
			//
			// TODO: implement https server startup
			//
			log.Fatalf("https server not implemented yet, exiting...")
		} else {
			rcv.server = &http.Server{
				Addr:           cli.String("server-address"),
				ReadTimeout:    cli.Duration("server-read-timeout"),
				WriteTimeout:   cli.Duration("server-write-timeout"),
				MaxHeaderBytes: int(cli.Int("server-max-header-bytes")),
				Handler:        r.Engine(),
			}

			if err := rcv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("http/s server stopped", "error", err.Error())
			}
		}
	}()
	return rcv
}

func (rcv *ServerProvider) Shutdown() error {
	cli := rcv.ctx.Value(common.KeyCommand).(*cli.Command)

	ctx, cancel := context.WithTimeout(
		context.Background(), cli.Duration("server-graceful-timeout"),
	)
	defer cancel()

	if err := rcv.server.Shutdown(ctx); err != nil {
		slog.Warn("http/s server shutted down", "error", err.Error())
		return err
	}
	slog.Info("http/s server shutted down", "error", "none")
	return nil
}
