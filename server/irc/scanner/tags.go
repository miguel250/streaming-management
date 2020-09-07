package scanner

import "github.com/miguel250/streaming-setup/server/irc/token"

func (sc *Scanner) scanTags(val *value, c rune) (token.Token, bool) {
	if c == '@' {
		sc.next()
		sc.startToken(val)
		sc.endToken(val)
		sc.scanningTag = true
		return token.AT, true
	}

	if sc.scanningTag {
		switch c {
		case ';':
			sc.next()
			sc.startToken(val)
			sc.endToken(val)
			return token.SEMICOLON, true
		case '=':
			sc.next()
			sc.startToken(val)
			return sc.scanTagValue(val), true
		default:
			sc.startToken(val)
			return sc.scanTagKey(val), true
		}
	}
	return 0, false
}

func (sc *Scanner) scanTagKey(val *value) token.Token {
	for {
		c, _ := sc.peak()

		if c == '=' {
			break
		}
		sc.next()
	}
	sc.endToken(val)
	return token.TAG
}

func (sc *Scanner) scanTagValue(val *value) token.Token {
	var c rune
	for {
		c, _ = sc.peak()

		if c == ';' || c == ' ' {
			break
		}

		sc.next()

	}
	sc.endToken(val)
	return token.EQUAL
}
