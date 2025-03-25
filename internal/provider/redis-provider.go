package provider

import (
	"context"

	"github.com/redis/go-redis/v9"

	"brickwall/internal/common"
)

type IRedisProvider interface {
	Connect() (*redis.Client, error)
	Client() *redis.Client
	Disconnect() error
}

type RedisProvider struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedisProvider(ctx context.Context) IRedisProvider {
	return &RedisProvider{ctx: ctx}

}

func (rcv *RedisProvider) Connect() (*redis.Client, error) {
	env := rcv.ctx.Value(common.KeyEnvProvider).(IEnvProvider)

	rcv.client = redis.NewClient(
		&redis.Options{
			Addr:       env.GetString("REDIS_ADDR", DefRedisAddr),
			Network:    env.GetString("REDIS_NETWORK", DefRedisNetwork),
			ClientName: env.GetString("REDIS_CLIENT_NAME", DefRedisClientName),
			DB:         env.GetInt("REDIS_DB", DefRedisDb),
			// TODO: impl other parameters
		},
	)
	if status := rcv.client.Ping(context.Background()); status.Err() != nil {
		rcv.client = nil
		return nil, status.Err()
	}
	return rcv.client, nil
}

func (rcv *RedisProvider) Client() *redis.Client {
	return rcv.client
}

func (rcv *RedisProvider) Disconnect() error {
	return rcv.client.Close()
}
