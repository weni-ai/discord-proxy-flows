package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DISCORD_PROXY_PORT", "9000")
	os.Setenv("DISCORD_PROXY_BOT_TOKEN", "my-bot-token")
	os.Setenv("DISCORD_PROXY_FLOWS_URL", "https://flows.ai")
	os.Setenv("DISCORD_PROXY_CHANNEL_UUID", "my-channel-uuid")

	config, err := LoadConfig()
	if err != nil {
		t.Errorf("Error loading config: %s", err.Error())
	}

	if config.Port != "9000" {
		t.Errorf("Expected Port to be '9000', but got '%s'", config.Port)
	}

	if config.BotToken != "my-bot-token" {
		t.Errorf("Expected BotToken to be 'my-bot-token', but got '%s'", config.BotToken)
	}

	if config.FlowsURL != "https://flows.ai" {
		t.Errorf("Expected FlowsURL to be 'https://flows.ai', but got '%s'", config.FlowsURL)
	}

	if config.ChannelUUID != "my-channel-uuid" {
		t.Errorf("Expected ChannelUUID to be 'my-channel-uuid', but got '%s'", config.ChannelUUID)
	}
}
