package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/miguel250/streaming-setup/server/irc"
	"github.com/miguel250/streaming-setup/server/twitch"
	twitch_util "github.com/miguel250/streaming-setup/server/twitch/util"
	emote_util "github.com/miguel250/streaming-setup/server/twitchemotes/util"
)

func CreateMockChatClient(t *testing.T) (*irc.Client, *EchoServer) {
	ts := MockTwitchChatServer(t)
	ts.Start()

	channeID := "558843277"
	testEndpoint := fmt.Sprintf("/kraken/users/%s", channeID)
	api, twitchMockServer := twitch_util.TestCreateClient(t, "user_response", testEndpoint, channeID)

	b, err := ioutil.ReadFile("testdata/channel_badges_response.json")

	if err != nil {
		t.Fatalf("Failed to open badges files with %s", err)
	}

	resp := &twitch.BadgesResponse{}

	err = json.Unmarshal(b, resp)

	if err != nil {
		t.Fatalf("Failed to parse json file with %s", err)
	}

	twitchEmotesMockAPI := emote_util.TestCreateClient(t, "emote_response", []string{"303365132"})

	conf := &irc.Config{
		Auth:         "test_auth_token",
		URL:          ts.addr,
		Name:         "test_account",
		Channel:      "test_channel",
		TwitchAPI:    api,
		TwitchEmotes: twitchEmotesMockAPI,
		Badges:       resp.BadgeSet,
	}

	client, err := irc.New(conf)

	if err != nil {
		t.Fatalf("Failed to connect to test server with %s", err)
	}

	t.Cleanup(func() {
		client.Close()
		ts.Shutdown()
		twitchMockServer.Close()
	})

	return client, ts
}
