package adapter

import (
	"context"

	"github.com/nats-io/nats.go"
)

type INatsAdapter interface {
	Publish(string, []byte) error
	Subscribe(string, func(*nats.Msg)) (*nats.Subscription, error)
}

type NatsAdapter struct {
	ctx  context.Context
	conn *nats.Conn
}

func NewNATSAdapter(ctx context.Context, conn *nats.Conn) INatsAdapter {
	return &NatsAdapter{ctx: ctx, conn: conn}
}

func (rcv *NatsAdapter) Publish(subject string, data []byte) error {
	return rcv.conn.Publish(subject, data)
}

func (rcv *NatsAdapter) Subscribe(subject string, handler func(*nats.Msg)) (*nats.Subscription, error) {
	return rcv.conn.Subscribe(subject, handler)
}
