package token

type Token int

const (
	EOF Token = iota
	INVALID

	COLON
	AT
	SEMICOLON
	EQUAL
	EXCLAMATION
	HASH
	WHITESPACE
	ASTERISK
	LEFTSQUAREBRACKET
	RIGHTSQUAREBRACKET

	// Command Responses
	WELCOME       // 001
	YOURHOST      // 002
	SERVERCREATED // 003
	SERVERMYINFO  // 004
	MOTDSTART     // 375
	MOTD          // 372
	MOTDEND       // 376
	NAMREPLY      // 353
	ENDOFNAMES    // 566

	// Commands
	CAP
	ACK
	PING
	JOIN
	PART
	CLEARCHAT
	CLEARMSG
	HOSTTARGET
	NOTICE
	RECONNECT
	ROOMSTATE
	USERNOTICE
	GLOBALUSERSTATE
	USERSTATE
	PRIVMSG

	// CAP
	TAG
)

func (t Token) String() string {
	return tokenName[t]
}

var tokenName = [...]string{
	EOF:                "eof",
	COLON:              "colon",
	AT:                 "at",
	SEMICOLON:          "semicolon",
	EQUAL:              "equal",
	EXCLAMATION:        "exclamation",
	HASH:               "hash",
	WHITESPACE:         "whitespace",
	INVALID:            "invalid token",
	ASTERISK:           "asterisk",
	LEFTSQUAREBRACKET:  "[",
	RIGHTSQUAREBRACKET: "]",

	WELCOME:       "welcome",
	YOURHOST:      "your host",
	SERVERCREATED: "server created",
	SERVERMYINFO:  "server my info",
	MOTDSTART:     "motdstart",
	MOTD:          "motd",
	MOTDEND:       "motdend",
	NAMREPLY:      "namreply",
	ENDOFNAMES:    "endofnames",

	CAP:             "cap",
	ACK:             "ack",
	PING:            "ping",
	JOIN:            "join",
	PART:            "part",
	CLEARCHAT:       "clear chat",
	CLEARMSG:        "clear msg",
	HOSTTARGET:      "host target",
	NOTICE:          "notice",
	RECONNECT:       "reconnect",
	ROOMSTATE:       "room state",
	USERNOTICE:      "user notice",
	GLOBALUSERSTATE: "global user state",
	USERSTATE:       "user state",
	PRIVMSG:         "private message",

	TAG: "tag",
}
