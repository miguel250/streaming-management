package parser_test

import (
	"fmt"
	"testing"

	"github.com/miguel250/streaming-setup/server/irc/parser"
	"github.com/miguel250/streaming-setup/server/irc/token"
)

func TestParseMsg(t *testing.T) {

	for _, test := range []struct {
		input    string
		channel  string
		message  string
		username string
		command  token.Token
		tags     map[string]string
	}{
		{
			`@badge-info=founder/2;badges=moderator/1,founder/0,bits/100;color=;display-name=AttackKopter;emotes=;flags=;id=a1d91e60-ae2b-4730-a2b1-38c23145887d;login=attackkopter;mod=1;msg-id=resub;msg-param-cumulative-months=2;msg-param-months=0;msg-param-should-share-streak=1;msg-param-streak-months=2;msg-param-sub-plan-name=Channel\sSubscription\s(miguelcodetv);msg-param-sub-plan=Prime;msg-param-was-gifted=false;room-id=558843277;subscriber=1;system-msg=AttackKopter\ssubscribed\swith\sTwitch\sPrime.\sThey've\ssubscribed\sfor\s2\smonths,\scurrently\son\sa\s2\smonth\sstreak!;tmi-sent-ts=1601065308944;user-id=239246205;user-type=mod :tmi.twitch.tv USERNOTICE #miguelcodetv :guess what`,
			"miguelcodetv",
			"guess what",
			"",
			token.USERNOTICE,
			map[string]string{
				"system-msg": "AttackKopter subscribed with Twitch Prime. They've subscribed for 2 months, currently on a 2 month streak!",
			},
		},
		{
			"@login=zpapa2112017;room-id=;target-msg-id=eec7a15c-ad91-45ac-a0ce-c52a2e8c9b65;tmi-sent-ts=1600803187681 :tmi.twitch.tv CLEARMSG #miguelcodetv :In search of followers, primes and views?",
			"miguelcodetv",
			"In search of followers, primes and views?",
			"",
			token.CLEARMSG,
			map[string]string{
				"login":         "zpapa2112017",
				"room-id":       "",
				"target-msg-id": "eec7a15c-ad91-45ac-a0ce-c52a2e8c9b65",
				"tmi-sent-ts":   "1600803187681",
			},
		},
		{
			":tmi.twitch.tv RECONNECT",
			"",
			"",
			"",
			token.RECONNECT,
			map[string]string{},
		},
		{
			":<user>!<user>@<user>.tmi.twitch.tv PRIVMSG #<channel> :This is a sample message",
			"<channel>",
			"This is a sample message",
			"<user>",
			token.PRIVMSG,
			map[string]string{},
		},
		{
			":tmi.twitch.tv CAP * ACK :twitch.tv/membership",
			"",
			"twitch.tv/membership",
			"",
			token.CAP,
			map[string]string{},
		},
		{
			":tmi.twitch.tv CAP * ACK :twitch.tv/tags",
			"",
			"twitch.tv/tags",
			"",
			token.CAP,
			map[string]string{},
		},
		{
			":ssp2014!ssp2014@ssp2014.tmi.twitch.tv JOIN #miguelcodetv",
			"miguelcodetv",
			"",
			"ssp2014",
			token.JOIN,
			map[string]string{},
		},
		{
			"PING :tmi.twitch.tv",
			"",
			"tmi.twitch.tv",
			"",
			token.PING,
			map[string]string{},
		},
		{
			":miguelcodetv_bot.tmi.twitch.tv 353 miguelcodetv_bot = #miguelcodetv :miguelcodetv_bot slaythor anotherttvviewer attackkopter miguelcodetv",
			"miguelcodetv",
			"miguelcodetv_bot slaythor anotherttvviewer attackkopter miguelcodetv",
			"",
			token.NAMREPLY,
			map[string]string{},
		},
		{
			"@badge-info=;badges=;color=;display-name=miguelcodetv_bot;emote-sets=0,564265402;user-id=567131665;user-type= :tmi.twitch.tv GLOBALUSERSTATE",
			"",
			"",
			"",
			token.GLOBALUSERSTATE,
			map[string]string{
				"badge-info":   "",
				"badges":       "",
				"color":        "",
				"display-name": "miguelcodetv_bot",
				"emote-sets":   "0,564265402",
				"user-id":      "567131665",
				"user-type":    "",
			},
		},
		{
			"@badge-info=;badges=;client-nonce=8ea2b6b2b091583b97d84454aefc6e2b;color=;display-name=sanjayshr;emotes=;flags=;id=63d172f6-a2f2-4d12-938b-a5be5b66a546;mod=0;room-id=558843277;subscriber=0;tmi-sent-ts=1598301071271;turbo=0;user-id=149222059;user-type= :sanjayshr!sanjayshr@sanjayshr.tmi.twitch.tv PRIVMSG #miguelcodetv :jwt ?",
			"miguelcodetv",
			"jwt ?",
			"sanjayshr",
			token.PRIVMSG,
			map[string]string{},
		},
		{
			"@badge-info=;badges=moderator/1;color=;display-name=miguelcodetv_bot;emote-sets=0,564265402;mod=1;subscriber=0;user-type=mod :tmi.twitch.tv USERSTATE #miguelcodetv",
			"miguelcodetv",
			"",
			"",
			token.USERSTATE,
			map[string]string{
				"badge-info":   "",
				"badges":       "moderator/1",
				"display-name": "miguelcodetv_bot",
				"emote-sets":   "0,564265402",
				"mod":          "1",
				"subscriber":   "0",
				"user-type":    "mod",
			},
		},
		{
			`@badge-info=;badges=premium/1;color=#008000;display-name=erikdotdev;emotes=;flags=;id=f1013215-e7e9-4441-830d-95bf7d12459f;login=erikdotdev;mod=0;msg-id=raid;msg-param-displayName=erikdotdev;msg-param-login=erikdotdev;msg-param-profileImageURL=https://static-cdn.jtvnw.net/jtv_user_pictures/2537a5a5-f45d-4cfb-80e2-f6b6b887ee23-profile_image-70x70.png;msg-param-viewerCount=44;room-id=558843277;subscriber=0;system-msg=44\sraiders\sfrom\serikdotdev\shave\sjoined!;tmi-sent-ts=1598300953914;user-id=192497221;user-type= :tmi.twitch.tv USERNOTICE #miguelcodetv`,
			"miguelcodetv",
			"",
			"",
			token.USERNOTICE,
			map[string]string{
				"badge-info":                "",
				"badges":                    "premium/1",
				"color":                     "#008000",
				"display-name":              "erikdotdev",
				"emotes":                    "",
				"flags":                     "",
				"id":                        "f1013215-e7e9-4441-830d-95bf7d12459f",
				"login":                     "erikdotdev",
				"mod":                       "0",
				"msg-id":                    "raid",
				"msg-param-displayName":     "erikdotdev",
				"msg-param-login":           "erikdotdev",
				"msg-param-profileImageURL": "https://static-cdn.jtvnw.net/jtv_user_pictures/2537a5a5-f45d-4cfb-80e2-f6b6b887ee23-profile_image-70x70.png",
				"msg-param-viewerCount":     "44",
				"room-id":                   "558843277",
				"subscriber":                "0",
				"system-msg":                `44 raiders from erikdotdev have joined!`,
				"tmi-sent-ts":               "1598300953914",
				"user-id":                   "192497221",
				"user-type":                 "",
			},
		},
	} {
		msg, err := parser.ParseMsg(test.input)

		if err != nil {
			fmt.Println(test.input)
			t.Fatalf("Failed to parse message with %s", err)
		}

		if msg.Channel != test.channel {
			t.Errorf("Channel doesn't match got: '%s', want: %s", msg.Channel, test.channel)
		}

		if msg.Message != test.message {
			t.Errorf("Message doesn't match got: '%s', want: %s", msg.Message, test.message)
		}

		if msg.Username != test.username {
			t.Errorf("Username doesn't match got: '%s', want: %s", msg.Username, test.username)
		}

		if msg.Command != test.command {
			t.Errorf("Command doesn't match got: '%s', want: %s", msg.Command, test.command)
		}

		for key, val := range test.tags {
			gotVal, ok := msg.Tags[key]

			if !ok {
				t.Errorf("Key is missing %s", key)
			}

			if val != gotVal {
				t.Errorf("Values didn't match want: %s, got: %s", val, gotVal)
			}
		}
	}
}
