package twitch_test

import (
	"testing"

	"github.com/miguel250/streaming-setup/server/cache"
	"github.com/miguel250/streaming-setup/server/twitch"
)

func TestConfigValidation(t *testing.T) {
	_, err := twitch.New(nil, cache.New())

	if err.Error() != twitch.ErrNilConf.Error() {
		t.Errorf("error didn't match what we expected %s", err)
	}

	_, err = twitch.New(&twitch.Config{}, cache.New())

	if err.Error() != twitch.ErrMissingTwitchURL.Error() {
		t.Errorf("error didn't match what we expected %s", err)
	}

	_, err = twitch.New(&twitch.Config{
		TwitchURL: "invalid_url",
	}, cache.New())

	if err.Error() != twitch.ErrMissingBadgeURL.Error() {
		t.Errorf("error didn't match what we expected %s", err)
	}

	_, err = twitch.New(&twitch.Config{
		TwitchURL: "invalid_url",
		BadgeURL:  "::////super_invalid_url",
	}, cache.New())

	if err.Error() != twitch.ErrMissingClientID.Error() {
		t.Errorf("error didn't match what we expected %s", err)
	}
}
