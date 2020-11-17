package scanner

import (
	"bytes"
	"testing"

	"github.com/miguel250/streaming-setup/server/irc/token"
)

func TestScanner(t *testing.T) {

	for _, test := range []struct {
		input, want string
	}{
		// issue #15
		{`:tmi.twitch.tv HOSTTARGET #hosting_channel :<channel> [<number-of-viewers>]`, "colon tmi.twitch.tv whitespace host target hash hosting_channel whitespace colon <channel> [<number-of-viewers>]"},
		// issue #10
		{`@badge-info=founder/2;badges=moderator/1,founder/0,bits/100;color=;display-name=AttackKopter;emotes=;flags=;id=a1d91e60-ae2b-4730-a2b1-38c23145887d;login=attackkopter;mod=1;msg-id=resub;msg-param-cumulative-months=2;msg-param-months=0;msg-param-should-share-streak=1;msg-param-streak-months=2;msg-param-sub-plan-name=Channel\sSubscription\s(miguelcodetv);msg-param-sub-plan=Prime;msg-param-was-gifted=false;room-id=558843277;subscriber=1;system-msg=AttackKopter\ssubscribed\swith\sTwitch\sPrime.\sThey've\ssubscribed\sfor\s2\smonths,\scurrently\son\sa\s2\smonth\sstreak!;tmi-sent-ts=1601065308944;user-id=239246205;user-type=mod :tmi.twitch.tv USERNOTICE #miguelcodetv :guess what`, `at tag badge-info equal founder/2 semicolon tag badges equal moderator/1,founder/0,bits/100 semicolon tag color equal semicolon tag display-name equal AttackKopter semicolon tag emotes equal semicolon tag flags equal semicolon tag id equal a1d91e60-ae2b-4730-a2b1-38c23145887d semicolon tag login equal attackkopter semicolon tag mod equal 1 semicolon tag msg-id equal resub semicolon tag msg-param-cumulative-months equal 2 semicolon tag msg-param-months equal 0 semicolon tag msg-param-should-share-streak equal 1 semicolon tag msg-param-streak-months equal 2 semicolon tag msg-param-sub-plan-name equal Channel\sSubscription\s(miguelcodetv) semicolon tag msg-param-sub-plan equal Prime semicolon tag msg-param-was-gifted equal false semicolon tag room-id equal 558843277 semicolon tag subscriber equal 1 semicolon tag system-msg equal AttackKopter\ssubscribed\swith\sTwitch\sPrime.\sThey've\ssubscribed\sfor\s2\smonths,\scurrently\son\sa\s2\smonth\sstreak! semicolon tag tmi-sent-ts equal 1601065308944 semicolon tag user-id equal 239246205 semicolon tag user-type equal mod whitespace colon tmi.twitch.tv whitespace user notice hash miguelcodetv whitespace colon guess what`},
		// issue: #8
		{"@login=zpapa2112017;room-id=;target-msg-id=eec7a15c-ad91-45ac-a0ce-c52a2e8c9b65;tmi-sent-ts=1600803187681 :tmi.twitch.tv CLEARMSG #miguelcodetv :In search of followers, primes and views?", "at tag login equal zpapa2112017 semicolon tag room-id equal semicolon tag target-msg-id equal eec7a15c-ad91-45ac-a0ce-c52a2e8c9b65 semicolon tag tmi-sent-ts equal 1600803187681 whitespace colon tmi.twitch.tv whitespace clear msg hash miguelcodetv whitespace colon In search of followers, primes and views?"},
		// issue: #9
		{":tmi.twitch.tv RECONNECT", "colon tmi.twitch.tv whitespace reconnect"},
		{":tmi.twitch.tv CAP * ACK :twitch.tv/membership", "colon tmi.twitch.tv whitespace cap whitespace asterisk whitespace ack whitespace colon twitch.tv/membership"},
		{":tmi.twitch.tv CAP * ACK :twitch.tv/tags", "colon tmi.twitch.tv whitespace cap whitespace asterisk whitespace ack whitespace colon twitch.tv/tags"},
		{":tmi.twitch.tv 001 <user> :Welcome, GLHF!", "colon tmi.twitch.tv whitespace welcome <user> whitespace colon Welcome, GLHF!"},
		{":tmi.twitch.tv 002 <user> :Your host is tmi.twitch.tv", "colon tmi.twitch.tv whitespace your host <user> whitespace colon Your host is tmi.twitch.tv"},
		{":tmi.twitch.tv 003 <user> :This server is rather new", "colon tmi.twitch.tv whitespace server created <user> whitespace colon This server is rather new"},
		{":tmi.twitch.tv 004 <user> :-", "colon tmi.twitch.tv whitespace server my info <user> whitespace colon -"},
		{":tmi.twitch.tv 375 <user> :-", "colon tmi.twitch.tv whitespace motdstart <user> whitespace colon -"},
		{":tmi.twitch.tv 372 <user> :You are in a maze of twisty passages.", "colon tmi.twitch.tv whitespace motd <user> whitespace colon You are in a maze of twisty passages."},
		{":tmi.twitch.tv 376 <user> :>", "colon tmi.twitch.tv whitespace motdend <user> whitespace colon >"},
		{"PING :tmi.twitch.tv", "ping whitespace colon tmi.twitch.tv"},
		{"PING :", "ping whitespace colon"},
		{":miguelcodetv_bot.tmi.twitch.tv 353 miguelcodetv_bot = #miguelcodetv :miguelcodetv_bot slaythor anotherttvviewer attackkopter miguelcodetv", "colon miguelcodetv_bot.tmi.twitch.tv whitespace namreply equal whitespace hash miguelcodetv whitespace colon miguelcodetv_bot slaythor anotherttvviewer attackkopter miguelcodetv"},
		{":miguelcodetv_bot.tmi.twitch.tv 366 miguelcodetv_bot #miguelcodetv :End of /NAMES list", "colon miguelcodetv_bot.tmi.twitch.tv whitespace endofnames hash miguelcodetv whitespace colon End of /NAMES list"},
		{":ronni!ronni@ronni.tmi.twitch.tv JOIN #dallas", "colon ronni exclamation ronni@ronni.tmi.twitch.tv whitespace join whitespace hash dallas"},
		{":ronni!ronni@ronni.tmi.twitch.tv PART #dallas", "colon ronni exclamation ronni@ronni.tmi.twitch.tv whitespace part hash dallas"},
		{":tmi.twitch.tv CLEARCHAT #dallas", "colon tmi.twitch.tv whitespace clear chat hash dallas"},
		{":tmi.twitch.tv CLEARCHAT #dallas :ronni", "colon tmi.twitch.tv whitespace clear chat hash dallas whitespace colon ronni"},
		{":tmi.twitch.tv CLEARMSG #dallas :HeyGuys", "colon tmi.twitch.tv whitespace clear msg hash dallas whitespace colon HeyGuys"},
		{":tmi.twitch.tv HOSTTARGET #hosting_channel :<channel> [1]", "colon tmi.twitch.tv whitespace host target hash hosting_channel whitespace colon <channel> [1]"},
		{":tmi.twitch.tv HOSTTARGET #hosting_channel :- [1]", "colon tmi.twitch.tv whitespace host target hash hosting_channel whitespace colon - [1]"},
		{":tmi.twitch.tv NOTICE #dallas :This room is no longer in slow mode.", "colon tmi.twitch.tv whitespace notice whitespace hash dallas whitespace colon This room is no longer in slow mode."},
		{":tmi.twitch.tv ROOMSTATE #<channel>", "colon tmi.twitch.tv whitespace room state whitespace hash <channel>"},
		{":tmi.twitch.tv USERNOTICE #<channel> :message", "colon tmi.twitch.tv whitespace user notice hash <channel> whitespace colon message"},
		{":tmi.twitch.tv USERSTATE #<channel>", "colon tmi.twitch.tv whitespace user state hash <channel>"},
		{":tmi.twitch.tv GLOBALUSERSTATE", "colon tmi.twitch.tv whitespace global user state"},
		{":<user>!<user>@<user>.tmi.twitch.tv PRIVMSG #<channel> :This is a sample message", "colon <user> exclamation <user>@<user>.tmi.twitch.tv whitespace private message hash <channel> whitespace colon This is a sample message"},
	} {
		scanner := NewScanner(test.input)
		var buf bytes.Buffer

		for {
			val, resultToken := scanner.NextToken()

			if resultToken == token.EOF {
				break
			}

			if resultToken == token.INVALID {
				t.Fatalf("Invalid token found %s on char '%s'", test.input, val.Text)
			}

			if buf.Len() > 0 {
				buf.WriteString(" ")
			}

			buf.WriteString(resultToken.String())
			if val.Text != "" {
				buf.WriteString(" ")
				buf.WriteString(val.Text)
			}
		}

		got := buf.String()

		if got != test.want {
			t.Errorf("Input didn't match output want: '%s', got: '%s'", test.want, got)
		}
	}
}

func TestTagScanner(t *testing.T) {

	for _, test := range []struct {
		input, want string
	}{
		{"@emote-only=0;followers-only=0;r9k=0;slow=0;subs-only=0 :tmi.twitch.tv ROOMSTATE #dallas", "at tag emote-only equal 0 semicolon tag followers-only equal 0 semicolon tag r9k equal 0 semicolon tag slow equal 0 semicolon tag subs-only equal 0 whitespace colon tmi.twitch.tv whitespace room state whitespace hash dallas"},
		{"@login=<login>;target-msg-id=<target-msg-id> :tmi.twitch.tv CLEARMSG #<channel> :<message>", "at tag login equal <login> semicolon tag target-msg-id equal <target-msg-id> whitespace colon tmi.twitch.tv whitespace clear msg hash <channel> whitespace colon <message>"},
		{"@login=ronni;target-msg-id=abc-123-def :tmi.twitch.tv CLEARMSG #dallas :HeyGuys", "at tag login equal ronni semicolon tag target-msg-id equal abc-123-def whitespace colon tmi.twitch.tv whitespace clear msg hash dallas whitespace colon HeyGuys"},
		{"@badge-info=<badge-info>;badges=<badges>;color=<color>;display-name=<display-name>;emote-sets=<emote-sets>;turbo=<turbo>;user-id=<user-id>;user-type=<user-type> :tmi.twitch.tv GLOBALUSERSTATE", "at tag badge-info equal <badge-info> semicolon tag badges equal <badges> semicolon tag color equal <color> semicolon tag display-name equal <display-name> semicolon tag emote-sets equal <emote-sets> semicolon tag turbo equal <turbo> semicolon tag user-id equal <user-id> semicolon tag user-type equal <user-type> whitespace colon tmi.twitch.tv whitespace global user state"},
		{"@badge-info=<badge-info>;badges=<badges>;color=<color>;display-name=<display-name>;emotes=<emotes>;id=<id-of-msg>;mod=<mod>;room-id=<room-id>;subscriber=<subscriber>;tmi-sent-ts=<timestamp>;turbo=<turbo>;user-id=<user-id>;user-type=<user-type> :<user>!<user>@<user>.tmi.twitch.tv PRIVMSG #<channel> :<message>", "at tag badge-info equal <badge-info> semicolon tag badges equal <badges> semicolon tag color equal <color> semicolon tag display-name equal <display-name> semicolon tag emotes equal <emotes> semicolon tag id equal <id-of-msg> semicolon tag mod equal <mod> semicolon tag room-id equal <room-id> semicolon tag subscriber equal <subscriber> semicolon tag tmi-sent-ts equal <timestamp> semicolon tag turbo equal <turbo> semicolon tag user-id equal <user-id> semicolon tag user-type equal <user-type> whitespace colon <user> exclamation <user>@<user>.tmi.twitch.tv whitespace private message hash <channel> whitespace colon <message>"},
		{"@emote-only=<emote-only>;followers-only=<followers-only>;r9k=<r9k>;slow=<slow>;subs-only=<subs-only> :tmi.twitch.tv ROOMSTATE #<channel>", "at tag emote-only equal <emote-only> semicolon tag followers-only equal <followers-only> semicolon tag r9k equal <r9k> semicolon tag slow equal <slow> semicolon tag subs-only equal <subs-only> whitespace colon tmi.twitch.tv whitespace room state whitespace hash <channel>"},
		{"@badge-info=<badge-info>;badges=<badges>;color=<color>;display-name=<display-name>;emotes=<emotes>;id=<id-of-msg>;login=<user>;mod=<mod>;msg-id=<msg-id>;room-id=<room-id>;subscriber=<subscriber>;system-msg=<system-msg>;tmi-sent-ts=<timestamp>;turbo=<turbo>;user-id=<user-id>;user-type=<user-type> :tmi.twitch.tv USERNOTICE #<channel> :<message>", "at tag badge-info equal <badge-info> semicolon tag badges equal <badges> semicolon tag color equal <color> semicolon tag display-name equal <display-name> semicolon tag emotes equal <emotes> semicolon tag id equal <id-of-msg> semicolon tag login equal <user> semicolon tag mod equal <mod> semicolon tag msg-id equal <msg-id> semicolon tag room-id equal <room-id> semicolon tag subscriber equal <subscriber> semicolon tag system-msg equal <system-msg> semicolon tag tmi-sent-ts equal <timestamp> semicolon tag turbo equal <turbo> semicolon tag user-id equal <user-id> semicolon tag user-type equal <user-type> whitespace colon tmi.twitch.tv whitespace user notice hash <channel> whitespace colon <message>"},
		{"@badge-info=<badge-info>;badges=<badges>;color=<color>;display-name=<display-name>;emote-sets=<emote-sets>;mod=<mod>;subscriber=<subscriber>;turbo=<turbo>;user-type=<user-type>:tmi.twitch.tv USERSTATE #<channel>", "at tag badge-info equal <badge-info> semicolon tag badges equal <badges> semicolon tag color equal <color> semicolon tag display-name equal <display-name> semicolon tag emote-sets equal <emote-sets> semicolon tag mod equal <mod> semicolon tag subscriber equal <subscriber> semicolon tag turbo equal <turbo> semicolon tag user-type equal <user-type>:tmi.twitch.tv whitespace user state hash <channel>"},
		{"@badge-info=;badges=;client-nonce=8ea2b6b2b091583b97d84454aefc6e2b;color=;display-name=sanjayshr;emotes=;flags=;id=63d172f6-a2f2-4d12-938b-a5be5b66a546;mod=0;room-id=558843277;subscriber=0;tmi-sent-ts=1598301071271;turbo=0;user-id=149222059;user-type= :sanjayshr!sanjayshr@sanjayshr.tmi.twitch.tv PRIVMSG #miguelcodetv :jwt ?", "at tag badge-info equal semicolon tag badges equal semicolon tag client-nonce equal 8ea2b6b2b091583b97d84454aefc6e2b semicolon tag color equal semicolon tag display-name equal sanjayshr semicolon tag emotes equal semicolon tag flags equal semicolon tag id equal 63d172f6-a2f2-4d12-938b-a5be5b66a546 semicolon tag mod equal 0 semicolon tag room-id equal 558843277 semicolon tag subscriber equal 0 semicolon tag tmi-sent-ts equal 1598301071271 semicolon tag turbo equal 0 semicolon tag user-id equal 149222059 semicolon tag user-type equal whitespace colon sanjayshr exclamation sanjayshr@sanjayshr.tmi.twitch.tv whitespace private message hash miguelcodetv whitespace colon jwt ?"},
	} {
		scanner := NewScanner(test.input)
		var buf bytes.Buffer

		for {
			val, resultToken := scanner.NextToken()

			if resultToken == token.EOF {
				break
			}

			if resultToken == token.INVALID {
				t.Fatalf("Invalid token found %s, invalid char: '%s'", test.input, val.Text)
			}

			if buf.Len() > 0 {
				buf.WriteString(" ")
			}

			buf.WriteString(resultToken.String())
			if val.Text != "" {
				buf.WriteString(" ")
				buf.WriteString(val.Text)
			}
		}

		got := buf.String()

		if got != test.want {
			t.Errorf("Input didn't match output want: '%s', got: '%s'", test.want, got)
		}
	}
}
