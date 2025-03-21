package provider

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urfave/cli/v3"

	"brickwall/internal/common"
)

type IPgxProvider interface {
	Open() (*pgxpool.Pool, error)
	Pool() *pgxpool.Pool
	Ping() error
	Close()
}

type PgxProvider struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewPgxProvider(ctx context.Context) IPgxProvider {
	return &PgxProvider{ctx: ctx}
}

func (rcv *PgxProvider) Open() (*pgxpool.Pool, error) {
	var (
		err  error
		conf *pgxpool.Config
	)
	//
	// TODO: implement SSL mode connection
	//
	cli := rcv.ctx.Value(common.KeyCommand).(*cli.Command)

	connUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cli.String("postgres-user"), cli.String("postgres-password"),
		cli.String("postgres-host"), cli.Int("postgres-port"), cli.String("postgres-db"),
	)
	if conf, err = pgxpool.ParseConfig(connUrl); err != nil {
		return nil, err
	}
	conf.MaxConns = int32(cli.Int("postgres-max-conns"))
	conf.MinConns = int32(cli.Int("postgres-min-conns"))
	conf.MaxConnLifetime = cli.Duration("postgres-max-conn-life-time")
	conf.MaxConnIdleTime = cli.Duration("postgres-max-conn-idle-time")
	conf.HealthCheckPeriod = cli.Duration("postgres-health-check-period")

	if rcv.pool, err = pgxpool.NewWithConfig(context.Background(), conf); err != nil {
		return nil, err
	}
	if err = rcv.pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return rcv.pool, nil
}

func (rcv *PgxProvider) Pool() *pgxpool.Pool {
	return rcv.pool
}

func (rcv *PgxProvider) Ping() error {
	return rcv.pool.Ping(context.Background())
}

func (rcv *PgxProvider) Close() {
	rcv.pool.Close()
}
