// Package config contains configuration utilities.
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config is a struct which contains the configuration for the application.
type Config struct {
	MongoDB struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"db"`
		URI      string `mapstructure:"uri"`
	} `mapstructure:"mongodb"`
	MQTT struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"mqtt"`
}

// Read reads in the configuration from the environment.
func Read() (*Config, error) {
	// MongoDB config options
	viper.SetDefault("mongodb.host", "")
	viper.SetDefault("mongodb.port", 27017)
	viper.SetDefault("mongodb.username", "")
	viper.SetDefault("mongodb.password", "")
	viper.SetDefault("mongodb.db", "home")
	viper.SetDefault("mongodb.uri", "")

	// MQTT config options
	viper.SetDefault("mqtt.host", "")
	viper.SetDefault("mqtt.port", 1883)
	viper.SetDefault("mqtt.username", "")
	viper.SetDefault("mqtt.password", "")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("BROKER")
	viper.AutomaticEnv()

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal configuration: %w", err)
	}

	return &config, nil
}
