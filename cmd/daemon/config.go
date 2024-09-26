package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     LoggerConfig     `mapstructure:"logger"`
	GrpcServer GrpcServerConfig `mapstructure:"grpcServer"`
	App        AppConfig        `mapstructure:"app"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type GrpcServerConfig struct {
	Port string `mapstructure:"port"`
}

type AppConfig struct {
	CollectCronSpec string        `mapstructure:"collectCronSpec"`
	CollectTimeout  time.Duration `mapstructure:"collectTimeout"`
	Metrics         MetricsConfig `mapstructure:"metrics"`
}

type MetricsConfig struct {
	La         bool `mapstructure:"la"`
	CPU        bool `mapstructure:"cpu"`
	DisksLoad  bool `mapstructure:"disksLoad"`
	Filesystem bool `mapstructure:"filesystem"`
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
