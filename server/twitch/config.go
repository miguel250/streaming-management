package twitch

import "errors"

var (
	ErrNilConf            = errors.New("twitch client config can't be nil")
	ErrMissingTwitchURL   = errors.New("twitch API url can't be empty")
	ErrMissingBadgeURL    = errors.New("twitch badge API url can't be empty")
	ErrMissingClientID    = errors.New("twitch client_id can't be empty")
	ErrMissingAuthURL     = errors.New("twitch auth url can't be empty")
	ErrMissingRedirectURL = errors.New("twitch redirect url can't be empty")
	ErrMissingSecret      = errors.New("twitch secret can't be empty")
)

type Config struct {
	TwitchURL   string
	ClientID    string
	BadgeURL    string
	Secret      string
	AuthURL     string
	RedirectURL string
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

	if conf.AuthURL == "" {
		return ErrMissingAuthURL
	}

	if conf.RedirectURL == "" {
		return ErrMissingRedirectURL
	}

	if conf.Secret == "" {
		return ErrMissingSecret
	}

	return nil
}
