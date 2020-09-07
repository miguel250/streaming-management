package scanner

import "github.com/miguel250/streaming-setup/server/irc/token"

func (sc *Scanner) scanWelcome(val *value, c rune) (token.Token, bool) {
	if c == '0' {
		c = sc.moveForward(3)

		switch c {
		case '1':
			sc.next()
			sc.startToken(val)
			c = sc.next()
			t := sc.scanSimpleMessage(val, c, token.WELCOME, false)
			sc.ignoreSpace = true
			return t, true
		case '2':
			sc.next()
			sc.startToken(val)
			c = sc.next()
			t := sc.scanSimpleMessage(val, c, token.YOURHOST, false)
			sc.ignoreSpace = true
			return t, true
		case '3':
			sc.next()
			sc.startToken(val)
			c = sc.next()
			t := sc.scanSimpleMessage(val, c, token.SERVERCREATED, false)
			sc.ignoreSpace = true
			return t, true
		case '4':
			sc.next()
			sc.startToken(val)
			c = sc.next()
			t := sc.scanSimpleMessage(val, c, token.SERVERMYINFO, false)
			sc.ignoreSpace = true
			return t, true
		}
	}
	return 0, false
}
