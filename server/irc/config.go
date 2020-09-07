package irc

import (
	"fmt"
	"reflect"

	"github.com/miguel250/streaming-setup/server/twitch"
	"github.com/miguel250/streaming-setup/server/twitchemotes"
)

type Config struct {
	Auth         string                          `json:"auth"`
	URL          string                          `json:"url"`
	Name         string                          `json:"name"`
	Channel      string                          `json:"channel"`
	TwitchAPI    *twitch.API                     `json:"-"`
	TwitchEmotes *twitchemotes.API               `json:"-"`
	Badges       map[string]*twitch.BadgeVersion `json:"-"`
}

func (c *Config) validate() error {

	if err := formatStrErr("auth", c.Auth); err != nil {
		return err
	}

	if err := formatStrErr("url", c.URL); err != nil {
		return err
	}

	if err := formatStrErr("name", c.Name); err != nil {
		return err
	}

	if err := formatStrErr("channel", c.Channel); err != nil {
		return err
	}

	if err := formatPtrErr("TwitchAPI", c.TwitchAPI); err != nil {
		return err
	}

	if err := formatPtrErr("TwitchEmotes", c.TwitchEmotes); err != nil {
		return err
	}

	return nil
}

func formatStrErr(fieldName string, value string) error {
	if value == "" {
		return fmt.Errorf("irc config: field %s can't be empty", fieldName)
	}
	return nil
}

func formatPtrErr(fieldName string, value interface{}) error {
	if reflect.ValueOf(value).IsNil() {
		return fmt.Errorf("irc config: field %s can't be empty", fieldName)
	}
	return nil
}
