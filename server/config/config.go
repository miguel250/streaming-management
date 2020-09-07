package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/miguel250/streaming-setup/server/irc"
)

type Config struct {
	path   string
	Twitch *Twitch `json:"twitch"`
}

type Twitch struct {
	FollowerGoalTotal int        `json:"follower_goal_total"`
	ClientID          string     `json:"client_id"`
	ChannelID         string     `json:"channel_id"`
	APIURL            string     `json:"api_url"`
	BadgesURL         string     `json:"badges_url"`
	IRC               irc.Config `json:"irc"`
	Emote             Emote      `json:"emote"`
}

type Emote struct {
	URL string `json:"url"`
}

func New(path string) (*Config, error) {
	body, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("failed to get configuration file with %w", err)
	}

	config := &Config{}
	err = json.Unmarshal(body, config)

	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration file with %w", err)
	}

	return config, nil
}
