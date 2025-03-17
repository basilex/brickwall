package provider

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v3"

	"brickwall/internal/common"
)

type IRedisProvider interface {
	Open() (*redis.Client, error)
	Client() *redis.Client
	Close() error
}

type RedisProvider struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedisProvider(ctx context.Context) IRedisProvider {
	return &RedisProvider{ctx: ctx}

}

func (rcv *RedisProvider) Open() (*redis.Client, error) {
	cli := rcv.ctx.Value(common.KeyCommand).(*cli.Command)

	rcv.client = redis.NewClient(
		&redis.Options{
			Addr:       cli.String("redis-addr"),
			Network:    cli.String("redis-network"),
			ClientName: cli.String("redis-client-name"),
			DB:         int(cli.Int("redis-db")),
			// TODO: impl other parameters
		},
	)
	status := rcv.client.Ping(context.Background())
	if status.Err() != nil {
		rcv.client = nil
		return nil, status.Err()
	}
	return rcv.client, nil
}

func (rcv *RedisProvider) Client() *redis.Client {
	return rcv.client
}

func (rcv *RedisProvider) Close() error {
	return rcv.client.Close()
}
