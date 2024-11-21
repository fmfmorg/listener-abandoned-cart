package config

import "os"

var (
	RedisAddress  = os.Getenv("FM_REDIS_ADDRESS")
	RedisPassword = os.Getenv("FM_REDIS_PASSWORD")

	ConnectChannel      = os.Getenv("FM_REDIS_WS_CONNECT_CHANNEL")
	DisconnectChannel   = os.Getenv("FM_REDIS_WS_DISCONNECT_CHANNEL")
	HeartbeatChannel    = os.Getenv("FM_REDIS_HEARTBEAT_CHANNEL")
	PaymentStartChannel = os.Getenv("FM_REDIS_PAYMENT_START_CHANNEL")
	PaymentEndChannel   = os.Getenv("FM_REDIS_PAYMENT_END_CHANNEL")

	RedisPaymentUpdateChannel = os.Getenv("FM_REDIS_PAYMENT_UPDATE_CHANNEL")

	HashKeyAllConnections = os.Getenv("FM_REDIS_HASHKEY_ALL_CONNECTIONS")

	WebshopApiUrl = os.Getenv("FM_CLIENT_WEBSHOP_API_URL")
)
