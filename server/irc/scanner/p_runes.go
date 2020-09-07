package scanner

import "github.com/miguel250/streaming-setup/server/irc/token"

func (sc *Scanner) scanRunesStartingP(val *value, c rune) (token.Token, bool) {
	if c == 'P' {
		c = sc.moveForward(2)

		switch c {
		case 'I':
			// PING command
			sc.moveForward(2)
			sc.startToken(val)
			sc.endToken(val)
			return token.PING, true
		case 'A':
			sc.moveForward(3)
			sc.startToken(val)
			sc.endToken(val)
			return token.PART, true
		case 'R':
			sc.moveForward(6)
			sc.startToken(val)
			sc.endToken(val)
			sc.ignoreSpace = true

			return token.PRIVMSG, true
		}
	}
	return 0, false
}
