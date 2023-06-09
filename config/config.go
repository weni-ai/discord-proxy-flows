package config

import (
	"github.com/jinzhu/configor"
)

type Config struct {
	Port        string `default:"8000" env:"DISCORD_PROXY_PORT"`
	BotToken    string `env:"DISCORD_PROXY_BOT_TOKEN"`
	FlowsURL    string `env:"DISCORD_PROXY_FLOWS_URL"`
	ChannelUUID string `env:"DISCORD_PROXY_CHANNEL_UUID"`
}

func LoadConfig() (Config, error) {
	var config Config

	settings := &configor.Config{
		ENVPrefix: "DISCORD_PROXY",
		Silent:    true,
	}

	err := configor.New(settings).Load(&config, "config.json")

	return config, err
}
