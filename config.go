package config

import (
	"io"
)

func New() *Config {
	return &Config{
		Bot: Bot{
			Commands: []string{},
		},
		Log: Log{
			Hooks: []LogHook{},
		},
		Database: Database{},
		Youtube:  Youtube{},
		Services: make(map[string]Service),
	}
}

// Config holds all config values
type Config struct {
	Bot      Bot                `mapstructure:"bot" json:"bot"`
	Log      Log                `mapstructure:"log" json:"log"`
	Database Database           `mapstructure:"database" json:"database"`
	Youtube  Youtube            `mapstructure:"youtube" json:"youtube"`
	Services map[string]Service `mapstructure:"services" json:"services"`
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
	Hooks []LogHook `mapstructure:"hooks" json:"hooks"`

	Level string `mapstructure:"level" json:"level"`
}

// LogHook represents a hook for the logger
type LogHook struct {
	Type     string `mapstructure:"type" json:"type"`
	Level    string `mapstructure:"level" json:"level"`
	URL      string `mapstructure:"url" json:"url"`
	Location string `mapstructure:"location" json:"location"`
}

// Database holds all config values for the database
type Database struct {
	Host     string `mapstructure:"host" json:"host"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	SSL      string `mapstructure:"ssl" json:"ssl"`
	Type     string `mapstructure:"type" json:"type"`
	Port     string `mapstructure:"port" json:"port"`
}

// Youtube holds all config values for youtube
type Youtube struct {
	APIKey string `mapstructure:"api-key" json:"api-key"`
}

// Service holds all config values for a service
type Service struct {
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
