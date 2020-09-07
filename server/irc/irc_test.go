package irc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/miguel250/streaming-setup/server/twitch"
	twitch_util "github.com/miguel250/streaming-setup/server/twitch/util"
	emote_util "github.com/miguel250/streaming-setup/server/twitchemotes/util"
)

func TestPrivMsgWithBadges(t *testing.T) {
	for _, test := range []struct {
		name         string
		input        string
		displayName  string
		channel      string
		message      string
		badges       []*twitch.Badge
		profileImage string
	}{
		{
			"testing tags",
			"@badge-info=;badges=;client-nonce=8ea2b6b2b091583b97d84454aefc6e2b;color=;display-name=sanjayshr;emotes=;flags=;id=63d172f6-a2f2-4d12-938b-a5be5b66a546;mod=0;room-id=558843277;subscriber=0;tmi-sent-ts=1598301071271;turbo=0;user-id=558843277;user-type= :sanjayshr!sanjayshr@sanjayshr.tmi.twitch.tv PRIVMSG #miguelcodetv :jwt ?",
			"sanjayshr",
			"miguelcodetv",
			"jwt ?",
			[]*twitch.Badge{},
			"https://static-cdn.jtvnw.net/jtv_user_pictures/cf98ab68-af25-441b-989e-f203cd46522e-profile_image-300x300.png",
		},
		{
			"testing badges",
			"@badge-info=founder/1;badges=moderator/1,founder/0,bits-leader/1,subscriber/0;client-nonce=2519a0bb9411a510293c39fac51323ad;color=;display-name=AttackKopter;emotes=;flags=;id=f12c675b-32b0-4ef3-8d20-e6c073ca6693;mod=1;room-id=558843277;subscriber=0;tmi-sent-ts=1599245324478;turbo=0;user-id=558843277;user-type=mod :attackkopter!attackkopter@attackkopter.tmi.twitch.tv PRIVMSG #miguelcodetv :wow",
			"AttackKopter",
			"miguelcodetv",
			"wow",
			[]*twitch.Badge{{
				Title:   "Subscriber",
				Image1X: "https://static-cdn.jtvnw.net/badges/v1/bea6cc27-c419-48e3-a121-110320d3482e/1",
				Image2X: "https://static-cdn.jtvnw.net/badges/v1/bea6cc27-c419-48e3-a121-110320d3482e/2",
				Image4X: "https://static-cdn.jtvnw.net/badges/v1/bea6cc27-c419-48e3-a121-110320d3482e/3",
			}},
			"https://static-cdn.jtvnw.net/jtv_user_pictures/cf98ab68-af25-441b-989e-f203cd46522e-profile_image-300x300.png",
		},
		{
			"testing emotes",
			"@badge-info=founder/1;badges=moderator/1,founder/0,bits-leader/1,subscriber/0;client-nonce=2519a0bb9411a510293c39fac51323ad;color=;display-name=AttackKopter;emotes=303365132;flags=;id=f12c675b-32b0-4ef3-8d20-e6c073ca6693;mod=1;room-id=558843277;subscriber=0;tmi-sent-ts=1599245324478;turbo=0;user-id=558843277;user-type=mod :attackkopter!attackkopter@attackkopter.tmi.twitch.tv PRIVMSG #miguelcodetv :wow miguel156Hero",
			"AttackKopter",
			"miguelcodetv",
			"wow <img src='https://static-cdn.jtvnw.net/emoticons/v1/303365132/2.0'>",
			[]*twitch.Badge{{
				Title:   "Subscriber",
				Image1X: "https://static-cdn.jtvnw.net/badges/v1/bea6cc27-c419-48e3-a121-110320d3482e/1",
				Image2X: "https://static-cdn.jtvnw.net/badges/v1/bea6cc27-c419-48e3-a121-110320d3482e/2",
				Image4X: "https://static-cdn.jtvnw.net/badges/v1/bea6cc27-c419-48e3-a121-110320d3482e/3",
			}},
			"https://static-cdn.jtvnw.net/jtv_user_pictures/cf98ab68-af25-441b-989e-f203cd46522e-profile_image-300x300.png",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			client, chatServerMock := createChatClient(t)
			var buf bytes.Buffer
			buf.WriteString(test.input)

			chatServerMock.setResponse(&buf)

			client.Start()
			data := <-client.OnMessage

			if data.DisplayName != test.displayName {
				t.Errorf("Displayname doesn't match want: %s, got: %s", test.displayName, data.DisplayName)
			}

			if data.Channel != test.channel {
				t.Errorf("Channel doesn't match want: %s, got: %s", test.channel, data.Channel)
			}

			if data.Message != test.message {
				t.Errorf("Message doesn't match want: %s, got: %s", test.message, data.Message)
			}

			if data.ProfileImage != test.profileImage {
				t.Errorf("Profile image url doesn't match want: %s, got: %s", test.profileImage, data.ProfileImage)
			}

			if len(data.Badges) != len(test.badges) {
				t.Errorf("Badges len to don't match got: %d, want: %d", len(data.Badges), len(test.badges))
			}

			if len(data.Badges) > 0 {
				for index, badge := range test.badges {
					gotBadge := data.Badges[index]

					if gotBadge.Title != badge.Title {
						t.Errorf("Badge tittle doesn't match got: %s, want: %s", gotBadge.Title, badge.Title)
					}

					if gotBadge.Image1X != badge.Image1X {
						t.Errorf("Image1X tittle doesn't match got: %s, want: %s", gotBadge.Image1X, badge.Image1X)
					}

					if gotBadge.Image2X != badge.Image2X {
						t.Errorf("Image2X tittle doesn't match got: %s, want: %s", gotBadge.Image2X, badge.Image2X)
					}

					if gotBadge.Image4X != badge.Image4X {
						t.Errorf("Image4X tittle doesn't match got: %s, want: %s", gotBadge.Image4X, badge.Image4X)
					}

				}
			}

		})
	}
}

func TestPrivMsg(t *testing.T) {
	client, chatServerMock := createChatClient(t)
	serverResponse := ":<user>!<user>@<user>.tmi.twitch.tv PRIVMSG #<channel> :This is a sample message"

	var buf bytes.Buffer
	buf.WriteString(serverResponse)

	chatServerMock.setResponse(&buf)

	client.Start()
	data := <-client.OnMessage

	wantDisplayname := "<user>"

	if data.DisplayName != wantDisplayname {
		t.Errorf("Displayname didn't match want: %s, got: %s", wantDisplayname, data.DisplayName)
	}

}

func TestSimpleCap(t *testing.T) {

	client, chatServerMock := createChatClient(t)
	serverResponse := ":tmi.twitch.tv CAP * ACK :twitch.tv/membership\n:tmi.twitch.tv CAP * ACK :twitch.tv/tags\n:tmi.twitch.tv CAP * ACK :twitch.tv/commands"

	var buf bytes.Buffer
	buf.WriteString(serverResponse)

	chatServerMock.setResponse(&buf)

	client.Start()
	err := client.Auth()

	if err != nil {
		t.Fatalf("Failed to auth with Twitch chat")
	}

	wants := []string{
		"twitch.tv/membership",
		"twitch.tv/tags",
		"twitch.tv/commands",
	}

	for _, want := range wants {
		got := <-client.OnCap

		if got.Message != want {
			t.Errorf("Commands don't match got: %s, want: %s", got.Message, want)
		}
	}
}

func createChatClient(t *testing.T) (*Client, *echoServer) {
	ts := testIRCEchoServer(t)
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

	conf := &Config{
		Auth:         "test_auth_token",
		URL:          ts.addr,
		Name:         "test_account",
		Channel:      "test_channel",
		TwitchAPI:    api,
		TwitchEmotes: twitchEmotesMockAPI,
		Badges:       resp.BadgeSet,
	}

	client, err := New(conf)

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
