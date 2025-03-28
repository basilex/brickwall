package provider

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

const (
	DefAppMode         string = "debug"
	DefEncoderStrategy string = "msgpack"

	DefTlsSslEnabled bool   = false
	DefTlsSslCert    string = "cert/server.crt"
	DefTlsSslKey     string = "cert/server.key"

	DefServerAddress         string        = "api:8081"
	DefServerReadTimeout     time.Duration = time.Duration(3 * time.Second)
	DefServerWriteTimeout    time.Duration = time.Duration(3 * time.Second)
	DefServerGracefulTimeout time.Duration = time.Duration(5 * time.Second)
	DefServerMaxHeaderBytes  int           = 1 << 20

	DefCorsEnabled          bool   = true
	DefCorsAllowOrigin      string = "*"
	DefCorsAllowHeaders     string = "Accept,Authorization,Content-Type,X-CSRF-Token"
	DefCorsAllowMethods     string = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	DefCorsExposeHeaders    string = "*"
	DefCorsAllowCredentials string = "false" // !no conv - cors related string
	DefCorsMaxAge           string = "300"   // !no conv - cors related string

	DefNatsURL                string        = nats.DefaultURL
	DefNatsMaxReconnect       int           = nats.DefaultMaxReconnect
	DefNatsReconnectWait      time.Duration = nats.DefaultReconnectWait
	DefNatsReconnectJitter    time.Duration = nats.DefaultReconnectJitter
	DefNatsReconnectJitterTLS time.Duration = nats.DefaultReconnectJitterTLS
	DefNatsTimeout            time.Duration = nats.DefaultTimeout
	DefNatsPingInterval       time.Duration = nats.DefaultPingInterval
	DefNatsMaxPingOut         int           = nats.DefaultMaxPingOut
	DefNatsReconnectBufSize   int           = nats.DefaultReconnectBufSize
	DefNatsDrainTimeout       time.Duration = nats.DefaultDrainTimeout
	DefNatsFlusherTimeout     time.Duration = nats.DefaultFlusherTimeout

	DefPostgresDb                string        = "bsp_dev"
	DefPostgresHost              string        = "postgres"
	DefPostgresPort              int           = 5432
	DefPostgresUser              string        = "system"
	DefPostgresPassword          string        = "passw0rd"
	DefPostgresMaxConns          int           = 100
	DefPostgresMinConns          int           = 5
	DefPostgresMaxConnLifeTime   time.Duration = time.Duration(10 * time.Minute)
	DefPostgresMaxConnIdleTime   time.Duration = time.Duration(3 * time.Minute)
	DefPostgresHealthCheckPeriod time.Duration = time.Duration(30 * time.Second)

	DefRedisAddr       string = "redis:6379"
	DefRedisNetwork    string = "tcp"
	DefRedisClientName string = "bsp"
	DefRedisDb         int    = 0

	DefJwtSecret            string        = "7b22fce240c32115056ba109f035542a3a1f9e1ee62fa653fa0a4ec0e267eb15"
	DefJwtAccessExpiration  time.Duration = time.Duration(15 * time.Minute)
	DefJwtRefreshExpiration time.Duration = time.Duration(24 * time.Hour)

	DefSmtpServerHost     string = "maildev"
	DefSmtpServerPort     int    = 1025
	DefSmtpServerUser     string = "system"
	DefSmtpServerPassword string = "passw0rd"
	DefSmtpSenderFrom     string = "no-replay@brickwall.com"
	DefSmtpTemplates      string = "templates"
	DefSmtpUseTLS         bool   = false
)

type EnvNS string

const (
	ConfigNS EnvNS = "config"
	SecretNS EnvNS = "secret"
)

type IEnvProvider interface {
	IsSwarm() bool

	GetString(string, string) string
	GetInt(string, int) int
	GetBool(string, bool) bool
	GetDuration(string, time.Duration) time.Duration

	Environment() map[string]string
}

type EnvProvider struct {
	ctx context.Context

	swarm  bool
	envMap map[string]string
}

func NewEnvProvider(ctx context.Context) IEnvProvider {
	envProvider := &EnvProvider{
		ctx:    ctx,
		swarm:  detectSwarm(),
		envMap: make(map[string]string),
	}
	if !envProvider.swarm {
		if err := godotenv.Load(); err != nil {
			slog.Error("bsp: failed to load <.env> environment")
			os.Exit(2)
		} else {
			slog.Info("bsp: environment <.env> loaded successfully")
		}
		for _, key := range os.Environ() {
			pair := strings.SplitN(key, "=", 2)
			if len(pair) == 2 {
				envProvider.envMap[pair[0]] = pair[1]
			}
		}
	}
	return envProvider
}

func detectSwarm() bool {
	if _, exists := os.LookupEnv("DOCKER_SERVICE_NAME"); exists {
		return true
	}
	if _, err := os.Stat("/var/run/" + string(SecretNS) + "/"); err == nil {
		return true
	}
	return false
}

func (rcv *EnvProvider) IsSwarm() bool {
	return rcv.swarm
}

func (rcv *EnvProvider) GetString(key string, defValue string) string {
	var ns EnvNS

	if rcv.secretData(key) {
		ns = SecretNS
	} else {
		ns = ConfigNS
	}
	if rcv.swarm {
		return rcv.lookupNsKey(ns, key)
	}
	if val, exists := rcv.envMap[key]; exists {
		return val
	}
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defValue
}

func (rcv *EnvProvider) GetInt(key string, defValue int) int {
	if val := rcv.GetString(key, ""); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		} else {
			slog.Error(
				"bsp", "message", "EnvProvider.GetInt() failure", "error", err,
			)
		}
	}
	return defValue
}

func (rcv *EnvProvider) GetBool(key string, defValue bool) bool {
	if val := rcv.GetString(key, ""); val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		} else {
			slog.Error(
				"bsp", "message", "EnvProvider.GetDuration() failure", "error", err,
			)
		}
	}
	return defValue
}

func (rcv *EnvProvider) GetDuration(key string, defValue time.Duration) time.Duration {
	if val := rcv.GetString(key, ""); val != "" {
		if durationVal, err := time.ParseDuration(val); err == nil {
			return durationVal
		} else {
			slog.Error(
				"bsp", "message", "EnvProvider.GetDuration() failure", "error", err,
			)
		}
	}
	return defValue
}

func (rcv *EnvProvider) Environment() map[string]string {
	newMap := make(map[string]string)

	for key, value := range rcv.envMap {
		if rcv.secretData(key) {
			newMap[key] = "***** hidden value *****"
			continue
		}
		newMap[key] = value
	}
	newMap["APP_SWARM_MODE"] = strconv.FormatBool(rcv.swarm)
	return newMap
}

// private functions
func (rcv *EnvProvider) lookupNsKey(ns EnvNS, key string) string {
	if rcv.swarm {
		path := "/run/" + string(ns) + "/" + key

		if data, err := os.ReadFile(path); err == nil {
			return strings.TrimSpace(string(data))
		}
	}
	if ns == ConfigNS {
		if val, exists := rcv.envMap[key]; exists {
			return val
		}
		if val, exists := os.LookupEnv(key); exists {
			return val
		}
	}
	slog.Warn(
		"bsp", string(ns), "ns", "key", key, "error", "failed to find candidate in ns/key",
	)
	return ""
}

func (rcv *EnvProvider) secretData(key string) bool {
	return strings.Contains(key, "SECRET") || strings.Contains(key, "PASSWORD")
}
