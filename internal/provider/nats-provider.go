package provider

// import (
// 	"brickwall/bsp/internal/common"
// 	"context"
// 	"errors"
// 	"log"

// 	"github.com/nats-io/nats.go"
// 	"github.com/urfave/cli/v3"
// )

// type INats interface {
// 	Connect() error
// 	Disconnect()
// 	Connection() *nats.Conn
// }

// type Nats struct {
// 	ctx  context.Context
// 	conn *nats.Conn
// }

// func NewNats(ctx context.Context) INats {
// 	return &Nats{ctx: ctx}
// }

// func (rcv *Nats) Connect() error {
// 	var err error

// 	cli := rcv.ctx.Value(common.KeyNats).(*cli.Command)

// 	connectionName := nats.Name(params.Name)
// 	connectTimeout := nats.Timeout(params.ConnectTimeout)
// 	reconnectWait := nats.ReconnectWait(params.ReconnectWait)
// 	pingInterval := nats.PingInterval(params.PingInterval)
// 	maxPingsOut := nats.MaxPingsOutstanding(params.MaxPingsOut)

// 	disconnectErrHandler := nats.DisconnectErrHandler(
// 		func(nc *nats.Conn, err error) {
// 			if !nc.IsClosed() {
// 				log.Printf("bus: disconnected due to: %s, will attempt reconnects for %.0fm", err, params.ConnectTimeout.Seconds())
// 			}
// 		})
// 	reconnectHandler := nats.ReconnectHandler(
// 		func(nc *nats.Conn) {
// 			log.Printf("bus: reconnected [%s]", nc.ConnectedUrl())
// 		})
// 	if rcv.conn, err = nats.Connect(
// 		params.ConnectURL,
// 		connectionName, connectTimeout, reconnectWait, pingInterval, maxPingsOut,
// 		disconnectErrHandler, reconnectHandler,
// 	); err != nil {
// 		return err
// 	}
// 	if rcv.conn.Status() != nats.CONNECTED {
// 		return errors.New("unable to establish failed")
// 	}
// 	return nil
// }

// func (rcv *Nats) Disconnect() {
// 	rcv.conn.Close()
// }

// func (rcv *Nats) Connection() *nats.Conn {
// 	return rcv.conn
// }
