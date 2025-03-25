package provider

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"brickwall/internal/common"
)

type IPgxProvider interface {
	Connect() (*pgxpool.Pool, error)
	Pool() *pgxpool.Pool
	Disconnect()
}

type PgxProvider struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewPgxProvider(ctx context.Context) IPgxProvider {
	return &PgxProvider{ctx: ctx}
}

func (rcv *PgxProvider) Connect() (*pgxpool.Pool, error) {
	var (
		err  error
		conf *pgxpool.Config
	)
	//
	// TODO: implement SSL mode connection
	//
	env := rcv.ctx.Value(common.KeyEnvProvider).(IEnvProvider)

	connUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		env.GetString("POSTGRES_USER", DefPostgresUser),
		env.GetString("POSTGRES_PASSWORD", DefPostgresPassword),
		env.GetString("POSTGRES_HOST", DefPostgresHost),
		env.GetInt("POSTGRES_PORT", DefPostgresPort),
		env.GetString("POSTGRES_DB", DefPostgresDb),
	)
	if conf, err = pgxpool.ParseConfig(connUrl); err != nil {
		return nil, err
	}
	conf.MaxConns = int32(env.GetInt("POSTGRES_MAX_CONNS", DefPostgresMaxConns))
	conf.MinConns = int32(env.GetInt("POSTGRES_MIN_CONNS", DefPostgresMinConns))
	conf.MaxConnLifetime = env.GetDuration("POSTGRES_MAX_CONN_LIFE_TIME", DefPostgresMaxConnLifeTime)
	conf.MaxConnIdleTime = env.GetDuration("POSTGRES_MAX_CONN_IDLE_TIME", DefPostgresMaxConnIdleTime)
	conf.HealthCheckPeriod = env.GetDuration("POSTGRES_HEALTH_CHECK_PERIOD", DefPostgresHealthCheckPeriod)

	if rcv.pool, err = pgxpool.NewWithConfig(rcv.ctx, conf); err != nil {
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

func (rcv *PgxProvider) Disconnect() {
	rcv.pool.Close()
}
