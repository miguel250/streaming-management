package irc

import (
	"fmt"
	"testing"

	"github.com/miguel250/streaming-setup/server/twitch"
)

func TestConfValidation(t *testing.T) {
	for _, test := range []struct {
		fieldName string
		conf      Config
	}{
		{
			"auth",
			Config{},
		},
		{
			"url",
			Config{
				Auth: "test",
			},
		},
		{
			"name",
			Config{
				Auth: "test",
				URL:  "http://example.com",
			},
		},
		{
			"channel",
			Config{
				Auth: "test",
				URL:  "http://example.com",
				Name: "account_name",
			},
		},
		{
			"TwitchAPI",
			Config{
				Auth:    "test",
				URL:     "http://example.com",
				Name:    "account_name",
				Channel: "test channel",
			},
		},
		{
			"TwitchEmotes",
			Config{
				Auth:      "test",
				URL:       "http://example.com",
				Name:      "account_name",
				Channel:   "test channel",
				TwitchAPI: &twitch.API{},
			},
		},
	} {
		t.Run(test.fieldName, func(t *testing.T) {
			err := test.conf.validate()

			if err == nil {
				t.Error("Expected error to not be nil")
			}

			want := fmt.Sprintf("irc config: field %s can't be empty", test.fieldName)
			if err.Error() != want {
				t.Errorf("Error message didn't match got: '%s', want: '%s'", err.Error(), want)
			}
		})
	}
}
