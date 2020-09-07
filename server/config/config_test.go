package config

import "testing"

func TestConfig(t *testing.T) {
	c, err := New("testdata/config.json")

	if err != nil {
		t.Fatalf("Failed test with %v", err)
	}

	var (
		wantClientID = "test_client_id"
		wantChannel  = "0001"
	)

	if c.Twitch.ClientID != wantClientID {
		t.Errorf("Client ID doesn't match got: %s, want: %s", c.Twitch.ClientID, wantClientID)
	}

	if c.Twitch.ChannelID != wantChannel {
		t.Errorf("Client ID doesn't match got: %s, want: %s", c.Twitch.ClientID, wantClientID)
	}
}

func TestMissingFile(t *testing.T) {
	_, err := New("invalid.json")

	if err == nil {
		t.Error("Should have returned an error")
	}

}
