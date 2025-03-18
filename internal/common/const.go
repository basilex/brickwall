package common

type KeyString string

const (
	KeyCommand  KeyString = "key-command"
	KeyMetadata KeyString = "key-metadata"

	KeyServiceManager KeyString = "key-service-manager"

	KeyValidatorProvider KeyString = "key-validator-provider"
	KeyRouterProvider    KeyString = "key-router-provider"
	KeyRedisProvider     KeyString = "key-redis-provider"
	KeyNatsProvider      KeyString = "key-nats-provider"
	KeyJwtProvider       KeyString = "key-jwt-provider"
	Key2FAProvider       KeyString = "key-2fa-provider"
	KeyPgxProvider       KeyString = "key-pgx-provider"
)
