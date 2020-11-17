package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/miguel250/streaming-setup/server/irc/scanner"
	"github.com/miguel250/streaming-setup/server/irc/token"
)

var (
	ErrEOF = errors.New("end of file")
)

type Message struct {
	raw          string
	scan         *scanner.Scanner
	Command      token.Token       `json:"-"`
	Username     string            `json:"username"`
	Channel      string            `json:"channel"`
	Message      string            `json:"message"`
	Tags         map[string]string `json:"tags"`
	currentToken token.Token
}

func (msg *Message) parseCap() error {
	if msg.currentToken == token.CAP {
		msg.Command = token.CAP

		for {
			val, err := msg.next()

			if err != nil {
				return err
			}

			if msg.currentToken == token.COLON {
				msg.Message = val
				return nil
			}
		}
	}
	return nil
}

func (msg *Message) parseTags() error {
	var currentKey string
	if msg.Command == 0 && msg.currentToken == token.AT {
		for {
			val, err := msg.next()
			if err != nil {
				return err
			}

			if msg.currentToken == token.TAG {
				msg.Tags[val] = ""
				currentKey = val
			}

			if currentKey != "" && msg.currentToken == token.EQUAL {
				// Tags can include multiple values marked by \s
				msg.Tags[currentKey] = strings.ReplaceAll(val, "\\s", " ")
			}

			if msg.currentToken == token.SEMICOLON {
				currentKey = ""
			}

			if msg.currentToken == token.WHITESPACE {
				return nil
			}
		}
	}
	return nil
}

func (msg *Message) parseUsername(val string) {
	if msg.Command == 0 && msg.currentToken == token.COLON {
		_, nexToken := msg.scan.NextToken()

		if nexToken == token.EXCLAMATION {
			msg.Username = val
		}
	}
}

func (msg *Message) parseJoin() error {
	if msg.currentToken == token.JOIN {
		msg.Command = token.JOIN

		for {
			val, err := msg.next()

			if err != nil {
				return err
			}

			if msg.currentToken == token.HASH {
				msg.Channel = val
				return nil
			}
		}
	}
	return nil
}
func (msg *Message) parsePing() error {
	if msg.currentToken == token.PING {
		msg.Command = token.PING
		for {
			val, err := msg.next()

			if err != nil {
				return err
			}

			if msg.currentToken == token.COLON {
				msg.Message = val
				return nil
			}
		}
	}
	return nil
}

func (msg *Message) parseNameReply() error {
	if msg.currentToken == token.NAMREPLY {
		msg.Command = token.NAMREPLY

		for {
			val, err := msg.next()

			if err != nil {
				return err
			}

			if msg.currentToken == token.HASH {
				msg.Channel = val
			}

			if msg.currentToken == token.COLON {
				msg.Message = val
			}

			if msg.Channel != "" && msg.Message != "" {
				return nil
			}
		}
	}
	return nil
}

func (msg *Message) parsePrivateMessage() error {
	return msg.parseMessage(token.PRIVMSG)
}

func (msg *Message) parseMessage(t token.Token) error {
	if msg.currentToken == t {
		msg.Command = t

		for {
			val, err := msg.next()
			if err != nil {
				return err
			}

			if msg.currentToken == token.HASH {
				msg.Channel = val
			}

			if msg.currentToken == token.COLON {
				msg.Message = val
			}

			if msg.Channel != "" && msg.Message != "" {
				break
			}
		}
	}
	return nil
}

func (msg *Message) parseSimpleCommandWithChannel(t token.Token) error {
	err := msg.parseSimpleCommands(t)
	if err != nil {
		return err
	}

	if msg.Command == t {
		for {
			val, err := msg.next()

			if err != nil {
				return nil
			}

			if msg.currentToken == token.HASH {
				msg.Channel = val
				return nil
			}
		}
	}
	return nil
}

func (msg *Message) parseSimpleCommands(token token.Token) error {
	if token == msg.currentToken {
		msg.Command = token
	}
	return nil
}

func (msg *Message) parseHostTarget(t token.Token) error {
	err := msg.parseSimpleCommands(t)
	if err != nil {
		return err
	}
	if msg.Command == token.HOSTTARGET {
		for {
			val, err := msg.next()
			if err != nil && err.Error() == ErrEOF.Error() {
				return nil
			}
			if err != nil {
				return err
			}

			if msg.currentToken == token.HASH {
				msg.Channel = val
			}

			if msg.currentToken == token.COLON {
				content := strings.Split(val, " ")
				size := len(content)
				if size <= 1 {
					return fmt.Errorf("failed to parse host target with invalid len (%d)", size)
				}
				msg.Tags["target_channel"] = content[0]
				msg.Message = content[1]
			}
		}
	}
	return err
}

func (msg *Message) next() (string, error) {
	if msg.currentToken != token.EOF {
		val, t := msg.scan.NextToken()

		if t == token.EOF {
			return "", ErrEOF
		}

		if t == token.INVALID {
			return "", fmt.Errorf("Failed to parse massage because of invalid token %s", msg.raw)
		}
		msg.currentToken = t
		return val.Text, nil
	}
	return "", ErrEOF
}

// ParseMsg parses a chat message sent from the Twitch chat server
// TODO: Add support for NOTICE, ROOMSTATE, HOSTTARGET, and CLEARCHAT
func ParseMsg(msg string) (*Message, error) {
	resultMsg := &Message{
		currentToken: -1,
		raw:          msg,
		scan:         scanner.NewScanner(msg),
		Tags:         make(map[string]string),
	}

	val, err := resultMsg.next()
	if err != nil {
		return nil, fmt.Errorf("failed to parse message with %w", err)
	}

	for {

		err := resultMsg.parseCap()
		if err != nil {
			return nil, err
		}

		err = resultMsg.parsePing()
		if err != nil {
			return nil, err
		}

		err = resultMsg.parseNameReply()
		if err != nil {
			return nil, err
		}

		err = resultMsg.parseTags()
		if err != nil {
			return nil, err
		}

		resultMsg.parseUsername(val)

		err = resultMsg.parseJoin()
		if err != nil {
			return nil, err
		}

		err = resultMsg.parsePrivateMessage()
		if err != nil {
			return nil, err
		}

		err = resultMsg.parseSimpleCommands(token.GLOBALUSERSTATE)
		if err != nil {
			return nil, err
		}

		err = resultMsg.parseSimpleCommandWithChannel(token.USERSTATE)
		if err != nil {
			return nil, err
		}

		err = resultMsg.parseHostTarget(token.HOSTTARGET)
		if err != nil {
			return nil, err
		}

		err = resultMsg.parseMessage(token.USERNOTICE)
		if err != nil {
			if err == ErrEOF {
				break
			}
			return nil, err
		}

		err = resultMsg.parseSimpleCommandWithChannel(token.RECONNECT)
		if err != nil {
			return nil, err
		}

		err = resultMsg.parseMessage(token.CLEARMSG)
		if err != nil {
			return nil, err
		}

		val, err = resultMsg.next()

		if err != nil && err == ErrEOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to parse message with %w", err)
		}
	}

	return resultMsg, nil
}
