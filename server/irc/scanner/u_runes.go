package scanner

import "github.com/miguel250/streaming-setup/server/irc/token"

func (sc *Scanner) scanRunesStartingU(val *value, c rune) (token.Token, bool) {

	if c == 'U' {
		c = sc.moveForward(5)

		if c == 'N' {
			sc.moveForward(6)
			sc.startToken(val)
			sc.endToken(val)
			sc.ignoreSpace = true
			return token.USERNOTICE, true
		}

		// USERSTATE
		if c == 'S' {
			sc.moveForward(5)
			sc.startToken(val)
			sc.endToken(val)
			return token.USERSTATE, true
		}
	}
	return 0, false
}
