package provider

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"

	"brickwall/internal/common"
)

type INatsProvider interface {
	Connect() (*nats.Conn, error)
	Connection() *nats.Conn
	JetStream() (nats.JetStreamContext, error)
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

	env := rcv.ctx.Value(common.KeyEnvProvider).(IEnvProvider)

	options := []nats.Option{
		nats.Name("BSP"),
		nats.MaxReconnects(env.GetInt("NATS_MAX_RECONNECT", DefNatsMaxReconnect)),
		nats.ReconnectWait(env.GetDuration("NATS_RECONNECT_WAIT", DefNatsReconnectWait)),
		nats.ReconnectJitter(
			env.GetDuration("NATS_RECONNECT_JITTER", DefNatsReconnectJitter),
			env.GetDuration("NATS_RECONNECT_JITTER_TLS", DefNatsReconnectJitterTLS),
		),
		nats.Timeout(env.GetDuration("NATS_TIMEOUT", DefNatsTimeout)),
		nats.PingInterval(env.GetDuration("NATS_PING_INTERVAL", DefNatsPingInterval)),
		nats.MaxPingsOutstanding(env.GetInt("NATS_MAX_PING_OUT", DefNatsMaxPingOut)),
		nats.ReconnectBufSize(int(env.GetInt("NATS_RECONNECT_BUF_SIZE", DefNatsReconnectBufSize))),
		nats.DrainTimeout(env.GetDuration("NATS_DRAIN_TIMEOUT", DefNatsDrainTimeout)),
		nats.FlusherTimeout(env.GetDuration("NATS_FLUSHER_TIMEOUT", DefNatsFlusherTimeout)),

		nats.DisconnectErrHandler(
			func(nc *nats.Conn, err error) {
				if !nc.IsClosed() {
					slog.Error(
						"nats", "error", err,
						"reconnects for", env.GetDuration("NATS_RECONNECT_WAIT", DefNatsReconnectWait),
					)
				}
			},
		),
		nats.ReconnectHandler(
			func(nc *nats.Conn) {
				slog.Warn("nats", "reconnected", nc.ConnectedUrl())
			},
		),
		nats.ErrorHandler(
			func(c *nats.Conn, s *nats.Subscription, err error) {
				slog.Error("nats", "error", err, "subscription", s.Subject)
			},
		),
		nats.ClosedHandler(
			func(nc *nats.Conn) {
				slog.Info("nats: connection closed")
			},
		),
	}
	natsURLs := env.GetString("NATS_URL", DefNatsURL)

	if rcv.conn, err = nats.Connect(natsURLs, options...); err != nil {
		return nil, err
	}
	return rcv.conn, nil
}

func (rcv *NatsProvider) Connection() *nats.Conn {
	return rcv.conn
}

func (rcv *NatsProvider) JetStream() (nats.JetStreamContext, error) {
	return rcv.conn.JetStream()
}

func (rcv *NatsProvider) Disconnect() {
	rcv.conn.Close()
}
