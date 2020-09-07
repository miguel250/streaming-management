package twitch

import "errors"

var (
	ErrNilConf          = errors.New("twitch client config can't be nil")
	ErrMissingTwitchURL = errors.New("twitch API url can't be empty")
	ErrMissingBadgeURL  = errors.New("twitch badge API url can't be empty")
	ErrMissingClientID  = errors.New("twitch client_id can't be empty")
)

type Config struct {
	TwitchURL string
	ClientID  string
	BadgeURL  string
}

func (conf *Config) validate() error {

	if conf.TwitchURL == "" {
		return ErrMissingTwitchURL
	}

	if conf.BadgeURL == "" {
		return ErrMissingBadgeURL
	}

	if conf.ClientID == "" {
		return ErrMissingClientID
	}
	return nil
}
