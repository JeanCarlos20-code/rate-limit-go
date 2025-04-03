package config

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type conf struct {
	jwtSecret    string `mapstructure:"JWT_SECRET"`
	jwtExpiresIn int    `mapstructure:"JWT_EXPIRESIN"`
	tokenAuth    *jwtauth.JWTAuth
}

func (c *conf) GetTokenAuth() *jwtauth.JWTAuth {
	return c.tokenAuth
}

func (c *conf) GetJWTExpiresIn() int {
	return c.jwtExpiresIn
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	cfg.jwtExpiresIn = viper.GetInt("JWT_EXPIRESIN")
	cfg.jwtSecret = viper.GetString("JWT_SECRET")

	cfg.tokenAuth = jwtauth.New("HS256", []byte(cfg.jwtSecret), nil)

	return cfg, err
}
