package scanner

import (
	"unicode/utf8"

	"github.com/miguel250/streaming-setup/server/irc/token"
)

type Scanner struct {
	input       []byte
	rest        []byte
	token       []byte
	pos         int
	ignoreSpace bool
	scanningTag bool
}

func NewScanner(msg string) *Scanner {
	b := []byte(msg)

	return &Scanner{
		input: b,
		rest:  b,
	}
}

type value struct {
	Text string
	pos  int
}

func (sc *Scanner) NextToken() (*value, token.Token) {
	var (
		c   rune
		val = &value{}
	)

	c, _ = sc.peak()

	if c == 0 {
		sc.startToken(val)
		sc.endToken(val)
		return val, token.EOF
	}

	if c == ' ' {
		sc.next()
		sc.startToken(val)
		sc.endToken(val)
		sc.scanningTag = false
		return val, token.WHITESPACE
	}

	scanToken, scanned := sc.scanTags(val, c)

	if scanned {
		return val, scanToken
	}

	switch c {
	case '[':
		sc.next()
		sc.startToken(val)
		return val, sc.scanSimpleMessage(val, token.LEFTSQUAREBRACKET, false)
	case ']':
		sc.next()
		sc.startToken(val)
		sc.endToken(val)
		return val, token.RIGHTSQUAREBRACKET
	case '!':
		// hostname
		sc.next()
		sc.startToken(val)
		return val, sc.scanSimpleMessage(val, token.EXCLAMATION, false)
	case '#':
		// Channel name
		sc.next()
		sc.startToken(val)
		return val, sc.scanSimpleMessage(val, token.HASH, false)
	case 'J':
		// JOIN
		sc.moveForward(4)
		sc.startToken(val)
		sc.endToken(val)
		return val, token.JOIN
	case ':':
		// Message value
		sc.next()
		sc.startToken(val)
		return val, sc.scanSimpleMessage(val, token.COLON, sc.ignoreSpace)
	case 'G':
		sc.moveForward(len("GLOBALUSERSTATE"))
		sc.startToken(val)
		sc.endToken(val)
		return val, token.GLOBALUSERSTATE
	case '*':
		sc.next()
		sc.startToken(val)
		sc.endToken(val)
		return val, token.ASTERISK
	case 'A':
		sc.moveForward(3)
		sc.startToken(val)
		sc.endToken(val)
		return val, token.ACK
	case 'H':
		// HOSTTARGET
		sc.moveForward(11)
		sc.startToken(val)
		sc.endToken(val)
		return val, token.HOSTTARGET
	case 'N':
		// NOTICE
		sc.moveForward(6)
		sc.startToken(val)
		sc.endToken(val)
		sc.ignoreSpace = true
		return val, token.NOTICE
	case '=':
		sc.next()
		sc.startToken(val)
		sc.endToken(val)
		return val, token.EQUAL
	}

	if scanToken, scanned := sc.scanRunesStartingC(val, c); scanned {
		return val, scanToken
	}

	// Welcome messages
	if scanToken, scanned := sc.scanWelcome(val, c); scanned {
		return val, scanToken
	}

	// General Responses codes
	if scanToken, scanned := sc.scanGeneralResponses(val, c); scanned {
		return val, scanToken
	}

	if scanToken, scanned := sc.scanRunesStartingP(val, c); scanned {
		return val, scanToken
	}

	if scanToken, scanned := sc.scanRunesStartingR(val, c); scanned {
		return val, scanToken
	}

	if scanToken, scanned := sc.scanRunesStartingU(val, c); scanned {
		return val, scanToken
	}

	val.Text = string(c)
	return val, token.INVALID
}

func (sc *Scanner) peak() (rune, int) {
	if len(sc.rest) == 0 {
		return 0, 0
	}

	return utf8.DecodeRune(sc.rest)
}

func (sc *Scanner) startToken(val *value) {
	sc.token = sc.rest
	val.pos = sc.pos
}

func (sc *Scanner) next() rune {
	if len(sc.rest) == 0 {
		panic("next at EOF")
	}

	c, size := sc.peak()
	sc.rest = sc.rest[size:]
	sc.pos++
	return c
}

func (sc *Scanner) moveForward(n int) rune {
	var c rune
	for i := 0; i < n; i++ {
		c = sc.next()
	}
	return c
}

func (sc *Scanner) scanSimpleMessage(val *value, t token.Token, ignoreSpace bool) token.Token {
	for {
		c, _ := sc.peak()

		if (!ignoreSpace && (c == ' ' || c == '!' || c == ']')) || c == 0 {
			break
		}
		sc.next()
	}
	sc.endToken(val)

	if ignoreSpace {
		sc.ignoreSpace = false
	}
	return t
}

func (sc *Scanner) endToken(val *value) {
	tokenSize := len(sc.token)
	restSize := len(sc.rest)
	textSize := tokenSize - restSize
	val.Text = string(sc.token[:textSize])
}
