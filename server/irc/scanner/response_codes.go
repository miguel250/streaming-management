package scanner

import "github.com/miguel250/streaming-setup/server/irc/token"

func (sc *Scanner) scanGeneralResponses(val *value, c rune) (token.Token, bool) {
	if c == '3' {
		c = sc.moveForward(2)

		// Message of the day codes
		if c == '7' {
			c = sc.next()
			switch c {
			case '5':
				sc.next()
				sc.startToken(val)
				sc.next()
				t := sc.scanSimpleMessage(val, token.MOTDSTART, false)
				sc.ignoreSpace = true
				return t, true
			case '2':
				sc.next()
				sc.startToken(val)
				sc.next()
				t := sc.scanSimpleMessage(val, token.MOTD, false)
				sc.ignoreSpace = true
				return t, true
			case '6':
				sc.next()
				sc.startToken(val)
				sc.next()
				t := sc.scanSimpleMessage(val, token.MOTDEND, false)
				sc.ignoreSpace = true
				return t, true
			}
		}

		if c == '5' {
			c = sc.next()
			if c == '3' {
				sc.next()
				sc.ignoreSpace = true
				return sc.scanNameReply(val)
			}
		}

		if c == '6' {
			c = sc.next()
			if c == '6' {
				sc.next()
				sc.ignoreSpace = true
				return sc.scanEndOfNames(val)
			}
		}
	}
	return 0, false
}

func (sc *Scanner) scanNameReply(val *value) (token.Token, bool) {
	for {
		c, _ := sc.peak()

		if c == '=' {
			break
		}
		sc.next()
	}

	sc.startToken(val)
	sc.endToken(val)
	return token.NAMREPLY, true
}

func (sc *Scanner) scanEndOfNames(val *value) (token.Token, bool) {
	sc.startToken(val)
	sc.endToken(val)

	for {
		c, _ := sc.peak()

		if c == '#' {
			break
		}
		sc.next()
	}
	return token.ENDOFNAMES, true
}
