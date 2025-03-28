package common

type KeyString string

const (
	// Metadata context keys
	KeyMetadata KeyString = "key-metadata"

	// Service manager context key
	KeyServiceManager KeyString = "key-service-manager"

	// Provider context keys
	KeyValidatorProvider KeyString = "key-validator-provider"
	KeyRouterProvider    KeyString = "key-router-provider"
	KeyRedisProvider     KeyString = "key-redis-provider"
	KeySmtpProvider      KeyString = "key-smtp-provider"
	KeyNatsProvider      KeyString = "key-nats-provider"
	KeyJwtProvider       KeyString = "key-jwt-provider"
	Key2FAProvider       KeyString = "key-2fa-provider"
	KeyPgxProvider       KeyString = "key-pgx-provider"
	KeyEnvProvider       KeyString = "key-env-provider"

	// Auth context keys
	KeyCtxUserID       KeyString = "key-ctx-user-id"
	KeyCtxAccessToken  KeyString = "key-ctx-access-token"
	KeyCtxRefreshToken KeyString = "key-ctx-refresh-token"
)

const (
	// NATS subjects/topics
	TopicUserRegistrationEmail KeyString = "topic.user.registration.email"
)

const (
	// JWT tokens
	JwtTokenValid   string = "valid"
	JwtTokenInvalid string = "invalid"
)
