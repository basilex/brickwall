package provider

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/urfave/cli/v3"

	"brickwall/internal/common"
)

type INatsProvider interface {
	Connect() (*nats.Conn, error)
	Connection() *nats.Conn
	Disconnect()
}

type NatsProvider struct {
	ctx  context.Context
	conn *nats.Conn
}

func NewNatsProvider(ctx context.Context) INatsProvider {
	return &NatsProvider{ctx: ctx}
}

func (rcv *NatsProvider) Connect() (*nats.Conn, error) {
	var err error

	cli := rcv.ctx.Value(common.KeyCommand).(*cli.Command)

	options := []nats.Option{
		nats.MaxReconnects(int(cli.Int("nats-max-reconnect"))),
		nats.ReconnectWait(cli.Duration("nats-reconnect-wait")),
		nats.ReconnectJitter(cli.Duration("nats-reconnect-jitter"), cli.Duration("nats-reconnect-jitter-tls")),
		nats.Timeout(cli.Duration("nats-timeout")),
		nats.PingInterval(cli.Duration("nats-ping-interval")),
		nats.MaxPingsOutstanding(int(cli.Int("nats-max-ping-out"))),
		nats.ReconnectBufSize(int(cli.Int("nats-reconnect-buf-size"))),
		nats.DrainTimeout(cli.Duration("nats-drain-timeout")),
		nats.FlusherTimeout(cli.Duration("nats-flusher-timeout")),

		nats.DisconnectErrHandler(
			func(nc *nats.Conn, err error) {
				if !nc.IsClosed() {
					slog.Error(
						"nats",
						"error", err,
						"reconnects for", cli.Duration("nats-reconnect-wait"),
					)
				}
			},
		),
		nats.ReconnectHandler(
			func(nc *nats.Conn) {
				slog.Warn(
					"nats",
					"reconnected", nc.ConnectedUrl(),
				)
			},
		),
		nats.ErrorHandler(
			func(c *nats.Conn, s *nats.Subscription, err error) {
				slog.Error(
					"nats",
					"error", err,
					"subscription", s.Subject,
				)
			},
		),
		nats.ClosedHandler(
			func(nc *nats.Conn) {
				slog.Info("nats: connection closed")
			},
		),
	}
	slog.Debug(">>>>>>>>>>>>>>", "nats-url", cli.String("nats-url"))
	if rcv.conn, err = nats.Connect(cli.String("nats-url"), options...); err != nil {
		return nil, err
	}
	return rcv.conn, nil
}

func (rcv *NatsProvider) Connection() *nats.Conn {
	return rcv.conn
}

func (rcv *NatsProvider) Disconnect() {
	rcv.conn.Close()
}
