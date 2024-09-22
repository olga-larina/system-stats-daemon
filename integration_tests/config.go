//go:build integration
// +build integration

package integration

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     LoggerConfig     `mapstructure:"logger"`
	GrpcClient GrpcClientConfig `mapstructure:"grpcClient"`
	App        AppConfig        `mapstructure:"app"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type GrpcClientConfig struct {
	GrpcURL string `mapstructure:"grpcUrl"`
}

type AppConfig struct {
	SendPeriodSeconds uint32 `mapstructure:"sendPeriodSeconds"`
	CalcPeriodSeconds uint32 `mapstructure:"calcPeriodSeconds"`
}

func NewConfig(path string) (*Config, error) {
	parser := viper.New()
	parser.SetConfigFile(path)

	err := parser.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	for _, key := range parser.AllKeys() {
		value := parser.GetString(key)
		parser.Set(key, os.ExpandEnv(value))
	}

	var config Config
	err = parser.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}

	return &config, err
}
