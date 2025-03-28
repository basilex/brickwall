package provider

import (
	"context"

	"github.com/nats-io/nats.go"

	"brickwall/internal/common"
)

const msgChanLen = 255

type MessageHandler func([]byte, IEncoder) error

type INatsSubscriber interface {
	Register(map[string]MessageHandler) error
	Shutdown()
}

type NatsSubscriber struct {
	nc       *nats.Conn
	ctx      context.Context
	subs     []*nats.Subscription
	cancel   context.CancelFunc
	encoder  IEncoder
	msgChans map[string]chan *nats.Msg
}

func NewNatsSubscriber(c context.Context) INatsSubscriber {
	envProvider := c.Value(common.KeyEnvProvider).(IEnvProvider)
	natsProvider := c.Value(common.KeyNatsProvider).(INatsProvider)

	encoder, _ := EncoderFactory(envProvider.GetString("ENCODER_STRATEGY", DefEncoderStrategy))

	ctx, cancel := context.WithCancel(c)

	return &NatsSubscriber{
		nc:       natsProvider.Connection(),
		ctx:      ctx,
		cancel:   cancel,
		encoder:  encoder,
		msgChans: make(map[string]chan *nats.Msg),
	}
}

func (rcv *NatsSubscriber) Register(topics map[string]MessageHandler) error {
	for topic, handler := range topics {
		ch := make(chan *nats.Msg, msgChanLen)
		rcv.msgChans[topic] = ch

		sub, err := rcv.nc.ChanSubscribe(topic, ch)
		if err != nil {
			return err
		}
		rcv.subs = append(rcv.subs, sub)
		go rcv.listen(topic, handler, rcv.encoder)
	}
	return nil
}

func (rcv *NatsSubscriber) listen(topic string, handler MessageHandler, encoder IEncoder) {
	for {
		select {
		case <-rcv.ctx.Done():
			return
		case msg := <-rcv.msgChans[topic]:
			if err := handler(msg.Data, encoder); err != nil {
			}
		}
	}
}

func (rcv *NatsSubscriber) Shutdown() {
	rcv.cancel()
	for _, sub := range rcv.subs {
		_ = sub.Unsubscribe()
	}
}
