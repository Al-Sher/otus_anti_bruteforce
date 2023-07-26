package config

import (
	"os"
	"strconv"
)

const (
	prefixEnv = "ANTI_BRUTEFORCE_"

	addrEnv          = "LISTEN_ADDR"
	loginLimitEnv    = "N"
	passwordLimitEnv = "M"
	ipLimitEnv       = "K"
	whiteListEnv     = "WHITE_LIST"
	blackListEnv     = "BLACK_LIST"
	redisURLEnv      = "REDIS_URL"
	logLevelEnv      = "LOG_LEVEL"
	bucketSizeEnv    = "BUCKET_SIZE"
	blockIntervalEnv = "BLOCK_INTERVAL"
	hostEnv          = "HOST"

	defaultAddr          = ":8081"
	defaultLoginLimit    = "10"
	defaultPasswordLimit = "100"
	defaultIPLimit       = "1000"
	defaultBucketSize    = "10"
	defaultBlockInterval = "3600"
	defaultWhiteListKey  = "ab:white"
	defaultBlackListKey  = "ab:black"
	defaultLogLevel      = "info"
	defaultRedisURL      = "redis://redis:6379/0"
	defaultHost          = "http://localhost" + defaultAddr
)

// Config структура конфигурации приложения.
type Config struct {
	Addr              string
	Host              string
	LoginLimit        int
	PasswordLimit     int
	IPLimit           int
	BucketSize        int
	BlockInterval     float64
	WhiteListRedisKey string
	BlackListRedisKey string
	LogLevel          string
	RedisURL          string
}

// New функция для создания конфигурации.
func New() (Config, error) {
	cfg := Config{}
	n, err := strconv.Atoi(Env(loginLimitEnv, defaultLoginLimit))
	if err != nil {
		return cfg, err
	}
	cfg.LoginLimit = n

	m, err := strconv.Atoi(Env(passwordLimitEnv, defaultPasswordLimit))
	if err != nil {
		return cfg, err
	}
	cfg.PasswordLimit = m

	k, err := strconv.Atoi(Env(ipLimitEnv, defaultIPLimit))
	if err != nil {
		return cfg, err
	}
	cfg.IPLimit = k

	bucketSize, err := strconv.Atoi(Env(bucketSizeEnv, defaultBucketSize))
	if err != nil {
		return cfg, err
	}
	cfg.BucketSize = bucketSize

	blockInterval, err := strconv.ParseFloat(Env(blockIntervalEnv, defaultBlockInterval), 64)
	if err != nil {
		return cfg, err
	}
	cfg.BlockInterval = blockInterval

	cfg.WhiteListRedisKey = Env(whiteListEnv, defaultWhiteListKey)
	cfg.BlackListRedisKey = Env(blackListEnv, defaultBlackListKey)
	cfg.LogLevel = Env(logLevelEnv, defaultLogLevel)
	cfg.RedisURL = Env(redisURLEnv, defaultRedisURL)
	cfg.Addr = Env(addrEnv, defaultAddr)
	cfg.Host = Env(hostEnv, defaultHost)

	return cfg, nil
}

func Env(key string, defaultValue string) string {
	v, ok := os.LookupEnv(prefixEnv + key)
	if !ok {
		return defaultValue
	}
	return v
}
