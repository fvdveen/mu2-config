package config

import (
	"io"
)

// Config holds all config values
type Config struct {
	Bot      Bot      `mapstructure:"bot" json:"bot"`
	Log      Log      `mapstructure:"log" json:"log"`
	Database Database `mapstructure:"database" json:"database"`
	Youtube  Youtube  `mapstructure:"youtube" json:"youtube"`
	Services Services `mapstructure:"services" json:"services"`
}

// Bot holds all config values for the bot
type Bot struct {
	Discord  Discord  `mapstructure:"discord" json:"discord"`
	Prefix   string   `mapstructure:"prefix" json:"prefix"`
	Commands []string `mapstructure:"commands" json:"commands"`
}

// Discord holds all config values for discord
type Discord struct {
	Token string `mapstructure:"token" json:"token"`
}

// Log holds all config values for the logger
type Log struct {
	Discord struct {
		Level   string `mapstructure:"level" json:"level"`
		WebHook string `mapstructure:"webhook" json:"webhook"`
	} `mapstructure:"discord" json:"discord"`

	Level string `mapstructure:"level" json:"level"`
}

// Database holds all config values for the database
type Database struct {
	Host     string `mapstructure:"host" json:"host"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	SSL      string `mapstructure:"ssl" json:"ssl"`
	Type     string `mapstructure:"type" json:"type"`
}

// Youtube holds all config values for youtube
type Youtube struct {
	APIKey string `mapstructure:"api-key" json:"api-key"`
}

// Services holds all config values for the services
type Services struct {
	Search Search `mapstructure:"location" json:"location"`
}

// Search holds all config values for the search service
type Search struct {
	Location string `mapstructure:"location" json:"location"`
}

// Watcher represents a stream of configs
type Watcher interface {
	Watch() <-chan *Config
}

// Provider represents a provider of configs
type Provider interface {
	Watcher
	io.Closer
}
