package api

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"github.com/urfave/cli/v3"

	"brickwall/cmd/api/service"
	"brickwall/internal/common"
	"brickwall/internal/provider"
)

var (
	defTlsSslEnabled bool   = false
	defTlsSslCert    string = "cert/server.crt"
	defTlsSslKey     string = "cert/server.key"

	defServerAddress         string        = "api:8081"
	defServerReadTimeout     time.Duration = time.Duration(3 * time.Second)
	defServerWriteTimeout    time.Duration = time.Duration(3 * time.Second)
	defServerGracefulTimeout time.Duration = time.Duration(5 * time.Second)
	defServerMaxHeaderBytes  int           = 1 << 20

	defCorsEnabled          bool   = true
	defCorsAllowOrigin      string = "*"
	defCorsAllowHeaders     string = "Accept,Authorization,Content-Type,X-CSRF-Token"
	defCorsAllowMethods     string = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	defCorsExposeHeaders    string = "*"
	defCorsAllowCredentials bool   = false
	defCorsMaxAge           int    = 300

	defNatsURL                string        = nats.DefaultURL
	defNatsMaxReconnect       int           = nats.DefaultMaxReconnect
	defNatsReconnectWait      time.Duration = nats.DefaultReconnectWait
	defNatsReconnectJitter    time.Duration = nats.DefaultReconnectJitter
	defNatsReconnectJitterTLS time.Duration = nats.DefaultReconnectJitterTLS
	defNatsTimeout            time.Duration = nats.DefaultTimeout
	defNatsPingInterval       time.Duration = nats.DefaultPingInterval
	defNatsMaxPingOut         int           = nats.DefaultMaxPingOut
	defNatsReconnectBufSize   int           = nats.DefaultReconnectBufSize
	defNatsDrainTimeout       time.Duration = nats.DefaultDrainTimeout
	defNatsFlusherTimeout     time.Duration = nats.DefaultFlusherTimeout

	defPostgresDb                string        = "bsp_dev"
	defPostgresHost              string        = "postgres"
	defPostgresPort              int           = 5432
	defPostgresUser              string        = "system"
	defPostgresPassword          string        = "passw0rd"
	defPostgresMaxConns          int           = 100
	defPostgresMinConns          int           = 5
	defPostgresMaxConnLifeTime   time.Duration = time.Duration(10 * time.Minute)
	defPostgresMaxConnIdleTime   time.Duration = time.Duration(3 * time.Minute)
	defPostgresHealthCheckPeriod time.Duration = time.Duration(30 * time.Second)

	defRedisAddr       string = "redis:6379"
	defRedisNetwork    string = "tcp"
	defRedisClientName string = "bsp"
	defRedisDb         int    = 0

	defJwtSecret            string        = "7b22fce240c32115056ba109f035542a3a1f9e1ee62fa653fa0a4ec0e267eb15"
	defJwtAccessExpiration  time.Duration = time.Duration(15 * time.Minute)
	defJwtRefreshExpiration time.Duration = time.Duration(24 * time.Hour)

	defSmtpServerHost     string = "maildev"
	defSmtpServerPort     int    = 1025
	defSmtpServerUser     string = "system"
	defSmtpServerPassword string = "passw0rd"
	defSmtpSenderFrom     string = "no-replay@brickwall.com"
	defSmtpTemplates      string = "templates"
	defSmtpUseTLS         bool   = false
)

func Command(ctx context.Context) *cli.Command {
	command := &cli.Command{
		Name:     "api",
		Category: "services",
		Usage:    "Run the api service",
		Action: func(ctx context.Context, cli *cli.Command) error {
			slog.Info("bsp: starting service api...")
			ctx = context.WithValue(ctx, common.KeyCommand, cli)
			return bootstrap(ctx)
		},
		Flags: []cli.Flag{
			//
			// TLS/SSL section
			//
			&cli.BoolFlag{
				Name:        "tls-ssl-enabled",
				Usage:       "Server HTTPS SSL enabled or not",
				Value:       defTlsSslEnabled,
				DefaultText: strconv.FormatBool(defTlsSslEnabled),
				Sources:     cli.EnvVars("TLS_SSL_ENABLED"),
			},
			&cli.StringFlag{
				Name:        "tls-ssl-cert",
				Usage:       "Server HTTPS SSL certificate cert",
				Value:       defTlsSslCert,
				DefaultText: defTlsSslCert,
				Sources:     cli.EnvVars("TLS_SSL_CERT"),
			},
			&cli.StringFlag{
				Name:        "tls-ssl-key",
				Usage:       "Server HTTPS SSL certificate key",
				Value:       defTlsSslKey,
				DefaultText: defTlsSslKey,
				Sources:     cli.EnvVars("TLS_SSL_KEY"),
			},
			//
			// Cors section
			//
			&cli.BoolFlag{
				Name:        "cors-enabled",
				Usage:       "Server HTTP CORS enabled or not",
				Value:       defCorsEnabled,
				DefaultText: strconv.FormatBool(defCorsEnabled),
				Sources:     cli.EnvVars("CORS_ENABLED"),
			},
			&cli.StringFlag{
				Name:        "cors-allow-origin",
				Usage:       "Server HTTP cors allowed origin",
				Value:       defCorsAllowOrigin,
				DefaultText: defCorsAllowOrigin,
				Sources:     cli.EnvVars("CORS_ALLOW_ORIGIN"),
			},
			&cli.StringFlag{
				Name:        "cors-allow-methods",
				Usage:       "Server HTTP cors allow methods",
				Value:       defCorsAllowMethods,
				DefaultText: defCorsAllowMethods,
				Sources:     cli.EnvVars("CORS_ALLOW_METHODS"),
			},
			&cli.StringFlag{
				Name:        "cors-allow-headers",
				Usage:       "Server HTTP cors allow headers",
				Value:       defCorsAllowHeaders,
				DefaultText: defCorsAllowHeaders,
				Sources:     cli.EnvVars("CORS_ALLOW_HEADERS"),
			},
			&cli.BoolFlag{
				Name:        "cors-allow-credentials",
				Usage:       "Server HTTP cors allow credentials",
				Value:       defCorsAllowCredentials,
				DefaultText: strconv.FormatBool(defCorsAllowCredentials),
				Sources:     cli.EnvVars("CORS_ALLOW_CREDENTIALS"),
			},
			&cli.StringFlag{
				Name:        "cors-expose-headers",
				Usage:       "Server HTTP cors expose headers",
				Value:       defCorsExposeHeaders,
				DefaultText: defCorsExposeHeaders,
				Sources:     cli.EnvVars("CORS_EXPOSE_HEADERS"),
			},
			&cli.IntFlag{
				Name:        "cors-max-age",
				Usage:       "Server max header bytes",
				Value:       int64(defCorsMaxAge),
				DefaultText: strconv.FormatInt(int64(defCorsMaxAge), 10),
				Sources:     cli.EnvVars("CORS_MAX_AGE"),
			},
			//
			// NATS section
			//
			&cli.StringFlag{
				Name:        "nats-url",
				Usage:       "NATS URL address",
				Value:       defNatsURL,
				DefaultText: defNatsURL,
				Sources:     cli.EnvVars("NATS_URL"),
			},
			&cli.IntFlag{
				Name:        "nats-max-reconnect",
				Usage:       "NATS max reconnects",
				Value:       int64(defNatsMaxReconnect),
				DefaultText: strconv.Itoa(defNatsMaxReconnect),
				Sources:     cli.EnvVars("NATS_MAX_RECONNECT"),
			},
			&cli.DurationFlag{
				Name:        "nats-reconnect-wait",
				Usage:       "NATS reconnect wait timeout",
				Value:       defNatsReconnectWait,
				DefaultText: defNatsReconnectWait.String(),
				Sources:     cli.EnvVars("NATS_RECONNECT_WAIT"),
			},
			&cli.DurationFlag{
				Name:        "nats-reconnect-jitter",
				Usage:       "NATS reconnect jitter timeout",
				Value:       defNatsReconnectJitter,
				DefaultText: defNatsReconnectJitter.String(),
				Sources:     cli.EnvVars("NATS_RECONNECT_JITTER"),
			},
			&cli.DurationFlag{
				Name:        "nats-reconnect-jitter-tls",
				Usage:       "NATS reconnect jitter tls timeout",
				Value:       defNatsReconnectJitterTLS,
				DefaultText: defNatsReconnectJitterTLS.String(),
				Sources:     cli.EnvVars("NATS_RECONNECT_JITTER_TLS"),
			},
			&cli.DurationFlag{
				Name:        "nats-timeout",
				Usage:       "NATS timeout",
				Value:       defNatsTimeout,
				DefaultText: defNatsTimeout.String(),
				Sources:     cli.EnvVars("NATS_TIMEOUT"),
			},
			&cli.DurationFlag{
				Name:        "nats-ping-interval",
				Usage:       "NATS ping interval",
				Value:       defNatsPingInterval,
				DefaultText: defNatsPingInterval.String(),
				Sources:     cli.EnvVars("NATS_PING_INTERVAL"),
			},
			&cli.IntFlag{
				Name:        "nats-max-ping-out",
				Usage:       "NATS max ping out",
				Value:       int64(defNatsMaxPingOut),
				DefaultText: strconv.Itoa(defNatsMaxPingOut),
				Sources:     cli.EnvVars("NATS_MAX_PING_OUT"),
			},
			&cli.IntFlag{
				Name:        "nats-reconnect-buf-size",
				Usage:       "NATS reconnect buffer size",
				Value:       int64(defNatsReconnectBufSize),
				DefaultText: strconv.Itoa(defNatsReconnectBufSize),
				Sources:     cli.EnvVars("NATS_RECONNECT_BUF_SIZE"),
			},
			&cli.DurationFlag{
				Name:        "nats-drain-timeout",
				Usage:       "NATS drain-timeout",
				Value:       defNatsDrainTimeout,
				DefaultText: defNatsDrainTimeout.String(),
				Sources:     cli.EnvVars("NATS_DRAIN_TIMEOUT"),
			},
			&cli.DurationFlag{
				Name:        "nats-flusher-timeout",
				Usage:       "NATS flusher-timeout",
				Value:       defNatsFlusherTimeout,
				DefaultText: defNatsFlusherTimeout.String(),
				Sources:     cli.EnvVars("NATS_FLUSHER_TIMEOUT"),
			},
			//
			// Server section
			//
			&cli.StringFlag{
				Name:        "server-address",
				Usage:       "Server HTTP binding address",
				Value:       defServerAddress,
				DefaultText: defServerAddress,
				Sources:     cli.EnvVars("SERVER_ADDRESS"),
			},
			&cli.DurationFlag{
				Name:        "server-read-timeout",
				Usage:       "Server read timeout",
				Value:       defServerReadTimeout,
				DefaultText: defServerReadTimeout.String(),
				Sources:     cli.EnvVars("SERVER_READ_TIMEOUT"),
			},
			&cli.DurationFlag{
				Name:        "server-write-timeout",
				Usage:       "Server write timeout",
				Value:       defServerWriteTimeout,
				DefaultText: defServerWriteTimeout.String(),
				Sources:     cli.EnvVars("SERVER_WRITE_TIMEOUT"),
			},
			&cli.DurationFlag{
				Name:        "server-graceful-timeout",
				Usage:       "Server shutdown graceful timeout",
				Value:       defServerGracefulTimeout,
				DefaultText: defServerGracefulTimeout.String(),
				Sources:     cli.EnvVars("SERVER_GRACEFUL_TIMEOUT"),
			},
			&cli.IntFlag{
				Name:        "server-max-header-bytes",
				Usage:       "Server max header bytes",
				Value:       int64(defServerMaxHeaderBytes),
				DefaultText: strconv.FormatInt(int64(defServerMaxHeaderBytes), 10),
				Sources:     cli.EnvVars("SERVER_MAX_HEADER_BYTES"),
			},
			//
			// Postgres section
			//
			&cli.StringFlag{
				Name:        "postgres-db",
				Usage:       "Postgres database name",
				Value:       defPostgresDb,
				DefaultText: defPostgresDb,
				Sources:     cli.EnvVars("POSTGRES_DB"),
			},
			&cli.StringFlag{
				Name:        "postgres-host",
				Usage:       "Postgres database host",
				Value:       defPostgresHost,
				DefaultText: defPostgresHost,
				Sources:     cli.EnvVars("POSTGRES_HOST"),
			},
			&cli.IntFlag{
				Name:        "postgres-port",
				Usage:       "Postgres database port",
				Value:       int64(defPostgresPort),
				DefaultText: strconv.FormatInt(int64(defPostgresPort), 10),
				Sources:     cli.EnvVars("POSTGRES_PORT"),
			},
			&cli.StringFlag{
				Name:        "postgres-user",
				Usage:       "Postgres database user",
				Value:       defPostgresUser,
				DefaultText: defPostgresUser,
				Sources:     cli.EnvVars("POSTGRES_USER"),
			},
			&cli.StringFlag{
				Name:        "postgres-password",
				Usage:       "Postgres database password",
				Value:       defPostgresPassword,
				DefaultText: defPostgresPassword,
				Sources:     cli.EnvVars("POSTGRES_PASSWORD"),
			},
			&cli.IntFlag{
				Name:        "postgres-max-conns",
				Usage:       "Postgres database max opened conns",
				Value:       int64(defPostgresMaxConns),
				DefaultText: strconv.FormatInt(int64(defPostgresMaxConns), 10),
				Sources:     cli.EnvVars("POSTGRES_MAX_CONNS"),
			},
			&cli.IntFlag{
				Name:        "postgres-min-conns",
				Usage:       "Postgres database min opened conns",
				Value:       int64(defPostgresMinConns),
				DefaultText: strconv.FormatInt(int64(defPostgresMinConns), 10),
				Sources:     cli.EnvVars("POSTGRES_MIN_CONNS"),
			},
			&cli.DurationFlag{
				Name:        "postgres-max-conn-life-time",
				Usage:       "Postgres database max conn life time",
				Value:       defPostgresMaxConnLifeTime,
				DefaultText: defPostgresMaxConnLifeTime.String(),
				Sources:     cli.EnvVars("POSTGRES_MAX_CONN_LIFE_TIME"),
			},
			&cli.DurationFlag{
				Name:        "postgres-max-conn-idle-time",
				Usage:       "Postgres database max conn idle time",
				Value:       defPostgresMaxConnIdleTime,
				DefaultText: defPostgresMaxConnIdleTime.String(),
				Sources:     cli.EnvVars("POSTGRES_MAX_CONN_IDLE_TIME"),
			},
			&cli.DurationFlag{
				Name:        "postgres-health-check-period",
				Usage:       "Postgres database health check period",
				Value:       defPostgresHealthCheckPeriod,
				DefaultText: defPostgresHealthCheckPeriod.String(),
				Sources:     cli.EnvVars("POSTGRES_HEALTH_CHECK_PERIOD"),
			},
			//
			// Redis section
			//
			&cli.StringFlag{
				Name:        "redis-addr",
				Usage:       "Redis server address",
				Value:       defRedisAddr,
				DefaultText: defRedisAddr,
				Sources:     cli.EnvVars("REDIS_ADDR"),
			},
			&cli.StringFlag{
				Name:        "redis-network",
				Usage:       "Redis network type (tcp|unix)",
				Value:       defRedisNetwork,
				DefaultText: defRedisNetwork,
				Sources:     cli.EnvVars("REDIS_NETWORK"),
			},
			&cli.StringFlag{
				Name:        "redis-client-name",
				Usage:       "Redis server client name",
				Value:       defRedisClientName,
				DefaultText: defRedisClientName,
				Sources:     cli.EnvVars("REDIS_CLIENT_NAME"),
			},
			&cli.IntFlag{
				Name:        "redis-db",
				Usage:       "Redis database number",
				Value:       int64(defRedisDb),
				DefaultText: strconv.FormatInt(int64(defRedisDb), 10),
				Sources:     cli.EnvVars("REDIS_DB"),
			},
			//
			// Jwt section
			//
			&cli.StringFlag{
				Name:        "jwt-secret",
				Usage:       "Jwt secret 32 bytes string",
				Value:       defJwtSecret,
				DefaultText: "Jwt secret string (32 bytes)",
				Sources:     cli.EnvVars("JWT_SECRET"),
			},
			&cli.DurationFlag{
				Name:        "jwt-access-expiration",
				Usage:       "Jwt access token expiration time",
				Value:       defJwtAccessExpiration,
				DefaultText: defJwtAccessExpiration.String(),
				Sources:     cli.EnvVars("JWT_ACCESS_EXPIRATION"),
			},
			&cli.DurationFlag{
				Name:        "jwt-refresh-expiration",
				Usage:       "Jwt refresh token expiration time",
				Value:       defJwtRefreshExpiration,
				DefaultText: defJwtRefreshExpiration.String(),
				Sources:     cli.EnvVars("JWT_REFRESH_EXPIRATION"),
			},
			//
			// SMTP settings
			//
			&cli.StringFlag{
				Name:        "smtp-server-host",
				Usage:       "SMTP server host",
				Value:       defSmtpServerHost,
				DefaultText: "SMTP server host address",
				Sources:     cli.EnvVars("SMTP_SERVER_HOST"),
			},
			&cli.IntFlag{
				Name:        "smtp-server-port",
				Usage:       "SMTP server port",
				Value:       int64(defSmtpServerPort),
				DefaultText: "SMTP server port number",
				Sources:     cli.EnvVars("SMTP_SERVER_PORT"),
			},
			&cli.StringFlag{
				Name:        "smtp-server-user",
				Usage:       "SMTP server user name",
				Value:       defSmtpServerUser,
				DefaultText: "SMTP server user name",
				Sources:     cli.EnvVars("SMTP_SERVER_USER"),
			},
			&cli.StringFlag{
				Name:        "smtp-server-password",
				Usage:       "SMTP server user password",
				Value:       defSmtpServerPassword,
				DefaultText: "SMTP server user password",
				Sources:     cli.EnvVars("SMTP_SERVER_PASSWORD"),
			},
			&cli.StringFlag{
				Name:        "smtp-sender-from",
				Usage:       "SMTP sender from address",
				Value:       defSmtpSenderFrom,
				DefaultText: "SMTP sender from adddress",
				Sources:     cli.EnvVars("SMTP_SENDER_FROM"),
			},
			&cli.StringFlag{
				Name:        "smtp-templates",
				Usage:       "SMTP templates dir",
				Value:       defSmtpTemplates,
				DefaultText: "SMTP templates dir",
				Sources:     cli.EnvVars("SMTP_TEMPLATES"),
			},
			&cli.BoolFlag{
				Name:        "smtp-use-tls",
				Usage:       "SMTP use TLS",
				Value:       defSmtpUseTLS,
				DefaultText: "SMTP use TLS",
				Sources:     cli.EnvVars("SMTP_USE_TLS"),
			},
		},
	}

	return command
}

// @title       Brickwall API
// @version     0.1.0
// @description This is Brickwall RestAPI
// @host        localhost:8081
// @BasePath    /api/v1
func bootstrap(ctx context.Context) error {
	//
	// Logger provider - no dependencies
	//
	slog.SetDefault(
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
	)
	//
	// NATS provider - no dependencies
	//
	natsProvider := provider.NewNatsProvider(ctx)
	if _, err := natsProvider.Connect(); err != nil {
		return err
	}
	ctx = context.WithValue(ctx, common.KeyNatsProvider, natsProvider)
	defer natsProvider.Disconnect()
	//
	// Redis provider - no dependencies
	//
	redisProvider := provider.NewRedisProvider(ctx)
	if _, err := redisProvider.Connect(); err != nil {
		return err
	}
	ctx = context.WithValue(ctx, common.KeyRedisProvider, redisProvider)
	defer redisProvider.Disconnect()
	//
	// Pgx provider - no dependencies
	//
	pgxProvider := provider.NewPgxProvider(ctx)
	if _, err := pgxProvider.Connect(); err != nil {
		return err
	}
	ctx = context.WithValue(ctx, common.KeyPgxProvider, pgxProvider)
	defer pgxProvider.Disconnect()
	//
	// Jwt provider - depends on Redis
	//
	jwtProvider := provider.NewJwtProvider(ctx)
	ctx = context.WithValue(ctx, common.KeyJwtProvider, jwtProvider)
	//
	// twoFA provider - no dependencies
	//
	twoFAProvider := provider.New2FAProvider(ctx)
	ctx = context.WithValue(ctx, common.Key2FAProvider, twoFAProvider)
	//
	// Router provider - no dependencies
	//
	routerProvider := provider.NewRouterProvider(ctx).Init()
	ctx = context.WithValue(ctx, common.KeyRouterProvider, routerProvider)
	//
	// Validator provider - no dependencies
	//
	validator := validator.New()
	ctx = context.WithValue(ctx, common.KeyValidatorProvider, validator)
	//
	// SMTP provider - no dependencies
	//
	smtpProvider := provider.NewSmtpProvider(ctx)
	ctx = context.WithValue(ctx, common.KeySmtpProvider, smtpProvider)
	//
	// Service Manager - depends on pgx and sqlc storage.queries
	//
	serviceManager := service.NewServiceManager(ctx)
	ctx = context.WithValue(ctx, common.KeyServiceManager, serviceManager)
	//
	// Server provider - no dependencies
	// Here is HTTP(S) server starts in goroutine
	//
	RegisterRoutes(ctx, routerProvider)
	srv := provider.NewServerProvider(ctx).Startup(routerProvider)
	defer srv.Shutdown()
	//
	// Watcher provider - no dependencies
	// Here the app is waiting for the signals to interruption
	//
	provider.NewWatcherProvider().Catch()
	return nil
}
