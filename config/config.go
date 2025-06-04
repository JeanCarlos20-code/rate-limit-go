package config

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type conf struct {
	jwtSecret     string `mapstructure:"JWT_SECRET"`
	jwtExpiresIn  int    `mapstructure:"JWT_EXPIRESIN"`
	tokenAuth     *jwtauth.JWTAuth
	redisAddr     string `mapstructure:"REDIS_ADDR"`
	redisPassword string `mapstructure:"REDIS_PASSWORD"`

	rateBackend   string `mapstructure:"RATE_BACKEND"`
	ipLimit       int    `mapstructure:"RATE_LIMIT_IP"`
	tokenLimit    int    `mapstructure:"RATE_LIMIT_TOKEN"`
	blockDuration int    `mapstructure:"BLOCK_DURATION_SECONDS"`
}

func (c *conf) GetRedisAddr() string {
	return c.redisAddr
}

func (c *conf) GetRedisPassword() string {
	return c.redisPassword
}

func (c *conf) GetTokenAuth() *jwtauth.JWTAuth {
	return c.tokenAuth
}

func (c *conf) GetJWTExpiresIn() int {
	return c.jwtExpiresIn
}

func (c *conf) GetRateBackend() string {
	return c.rateBackend
}

func (c *conf) GetIPLimit() int {
	return c.ipLimit
}

func (c *conf) GetTokenLimit() int {
	return c.tokenLimit
}

func (c *conf) GetBlockDuration() int {
	return c.blockDuration
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	cfg.jwtSecret = viper.GetString("JWT_SECRET")
	cfg.jwtExpiresIn = viper.GetInt("JWT_EXPIRESIN")

	cfg.redisAddr = viper.GetString("REDIS_ADDR")
	cfg.redisPassword = viper.GetString("REDIS_PASSWORD")

	cfg.rateBackend = viper.GetString("RATE_BACKEND")
	cfg.ipLimit = viper.GetInt("RATE_LIMIT_IP")
	cfg.tokenLimit = viper.GetInt("RATE_LIMIT_TOKEN")
	cfg.blockDuration = viper.GetInt("BLOCK_DURATION_SECONDS")

	cfg.tokenAuth = jwtauth.New("HS256", []byte(cfg.jwtSecret), nil)

	return cfg, nil
}
