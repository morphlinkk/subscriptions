package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort              string        `mapstructure:"SERVER_PORT"`
	EnvMode                 string        `mapstructure:"ENV_MODE"`
	LogLevel                string        `mapstructure:"LOG_LEVEL"`
	DatabaseURI             string        `mapstructure:"DATABASE_URI"`
	DatabaseMaxConnections  int           `mapstructure:"DATABASE_MAXCONNS"`
	DatabaseMinConnections  int           `mapstructure:"DATABASE_MINCONNS"`
	DatabaseMaxConnLifetime time.Duration `mapstructure:"DATABASE_MAXCONNLIFETIME"`
}

func Load() (*Config, error) {
	v := viper.New()
	setDefaults(v)

	v.SetConfigFile(".env")
	v.SetConfigType("env")

	var fileLookupError viper.ConfigFileNotFoundError

	if err := v.ReadInConfig(); err != nil {
		if errors.As(err, &fileLookupError) {
			return nil, fmt.Errorf("unable to find config file")
		} else {
			return nil, fmt.Errorf("unable to read in config file")
		}
	}

	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("SERVER_PORT", "3000")
	v.SetDefault("ENV_MODE", "debug")
	v.SetDefault("LOG_LEVEL", "debug")
	v.SetDefault("DATABASE_MAXCONNECTIONS", 25)
	v.SetDefault("DATABASE_MINCONNECTIONS", 5)
	v.SetDefault("DATABASE_MAXCONNLIFETIME", 30*time.Minute)
}

func (c *Config) Validate() error {
	if c.DatabaseURI == "" {
		return fmt.Errorf("DATABASE_URI is required but not set")
	}

	if c.DatabaseMaxConnections < c.DatabaseMinConnections {
		return fmt.Errorf("DATABASE_MAXCONNECTIONS must be >= DATABASE_MINCONNECTIONS")
	}

	if c.DatabaseMaxConnections < 1 {
		return fmt.Errorf("DATABASE_MAXCONNECTIONS must be at least 1")
	}

	return nil
}
