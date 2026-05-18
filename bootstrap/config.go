package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Env string `mapstructure:"env"`
}

type ServerConfig struct {
	Address        string `mapstructure:"address"`
	ContextTimeout int    `mapstructure:"context_timeout"`
	PublicURL      string `mapstructure:"public_url"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type JWTConfig struct {
	AccessTokenExpiryHour  int    `mapstructure:"access_token_expiry_hour"`
	RefreshTokenExpiryHour int    `mapstructure:"refresh_token_expiry_hour"`
	AccessTokenSecret      string `mapstructure:"access_token_secret"`
	RefreshTokenSecret     string `mapstructure:"refresh_token_secret"`
}

type MailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Mail     MailConfig     `mapstructure:"mail"`
}

func NewConfig() *Config {
	v := viper.New()
	v.SetConfigFile("config.yaml")

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("Can't load config.yaml: ", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatal("Config can't be loaded: ", err)
	}

	if cfg.App.Env == "development" {
		log.Println("The App is running in development env")
	}

	return &cfg
}
