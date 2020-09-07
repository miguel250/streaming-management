package scanner

import "github.com/miguel250/streaming-setup/server/irc/token"

func (sc *Scanner) scanRunesStartingR(val *value, c rune) (token.Token, bool) {
	if c == 'R' {
		c = sc.moveForward(2)
		if c == 'O' {
			// ROOMSTATE
			sc.moveForward(7)
			sc.startToken(val)
			sc.endToken(val)
			return token.ROOMSTATE, true
		}
	}

	return 0, false
}
