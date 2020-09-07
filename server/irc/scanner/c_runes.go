package scanner

import "github.com/miguel250/streaming-setup/server/irc/token"

func (sc *Scanner) scanRunesStartingC(val *value, c rune) (token.Token, bool) {
	if c == 'C' {
		sc.next()
		c = sc.next()

		switch c {
		case 'A':
			//CAP
			sc.next()
			sc.startToken(val)
			sc.endToken(val)
			return token.CAP, true
		case 'L':
			c = sc.moveForward(4)
			switch c {
			case 'C':
				sc.moveForward(4)
				sc.startToken(val)
				sc.endToken(val)
				return token.CLEARCHAT, true
			case 'M':
				sc.moveForward(3)
				sc.startToken(val)
				sc.endToken(val)
				return token.CLEARMSG, true
			}
		}
	}
	return 0, false
}
