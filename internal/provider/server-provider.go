package provider

import (
	"context"
	"log"
	"log/slog"
	"net/http"

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
	env := rcv.ctx.Value(common.KeyEnvProvider).(IEnvProvider)

	go func() {
		slog.Info(
			"http/s server started", "bind", env.GetString("SERVER_ADDRESS", DefServerAddress),
		)
		if env.GetBool("TLS_SSL_ENABLED", DefTlsSslEnabled) {
			//
			// TODO: implement https server startup
			//
			log.Fatalf("https server not implemented yet, exiting...")
		} else {
			rcv.server = &http.Server{
				Addr:           env.GetString("SERVER_ADDRESS", DefServerAddress),
				ReadTimeout:    env.GetDuration("SERVER_READ_TIMEOUT", DefServerReadTimeout),
				WriteTimeout:   env.GetDuration("SERVER_WRITE_TIMEOUT", DefServerWriteTimeout),
				MaxHeaderBytes: int(env.GetInt("SERVER_MAX_HEADER_BYTES", DefServerMaxHeaderBytes)),
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
	env := rcv.ctx.Value(common.KeyEnvProvider).(IEnvProvider)

	ctx, cancel := context.WithTimeout(
		context.Background(), env.GetDuration("SERVER_GRACEFUL_TIMEOUT", DefServerGracefulTimeout),
	)
	defer cancel()

	if err := rcv.server.Shutdown(ctx); err != nil {
		slog.Warn("http/s server shutted down", "error", err.Error())
		return err
	}
	slog.Info("http/s server shutted down", "error", "none")
	return nil
}
