package irc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"sync"

	"github.com/miguel250/streaming-setup/server/irc/parser"
	"github.com/miguel250/streaming-setup/server/irc/token"
	"github.com/miguel250/streaming-setup/server/twitch"
	"github.com/miguel250/streaming-setup/server/twitchemotes"
)

type chatCommand int

const (
	Cap chatCommand = iota
	Pass
	Nick
	Join
	Pong
	PrivMsg
)

var commandToString = map[chatCommand]string{
	Cap:     "CAP REQ :",
	Pass:    "PASS oauth:",
	Nick:    "NICK ",
	Join:    "JOIN #",
	Pong:    "PONG :",
	PrivMsg: "PRIVMSG #",
}

type Client struct {
	sync.RWMutex
	conf         *Config
	conn         net.Conn
	close        chan bool
	OnCap        chan *parser.Message
	OnMessage    chan *Message
	twitchEmotes *twitchemotes.API
	twitchClient *twitch.API
	badges       map[string]*twitch.BadgeVersion
	currentUsers map[string]*user
	emotesCache  map[string]*emote
}

type Message struct {
	Badges       []*twitch.Badge `json:"badges"`
	DisplayName  string          `json:"display-name"`
	Message      string          `json:"message"`
	ProfileImage string          `json:"profile_image"`
	Channel      string          `json:"channel"`
}

type user struct {
	displayName  string
	profileImage string
}

type emote struct {
	id       string
	code     string
	imageURL string
}

func (c *Client) Start() error {
	conn, err := net.Dial("tcp", c.conf.URL)

	if err != nil {
		return fmt.Errorf("Failed to create connection with %w", err)
	}

	c.conn = conn

	reader := bufio.NewReader(c.conn)
	tp := textproto.NewReader(reader)

	go func() {
		for {
			line, err := tp.ReadLine()

			if err != nil && err == io.EOF {
				fmt.Println("Connection was closed")
				return
			}

			if err != nil {
				return
			}

			parse, err := parser.ParseMsg(line)

			if err != nil {
				fmt.Printf("Failed to parse message with %s\n", err)
			}

			switch parse.Command {
			case token.CAP:
				c.OnCap <- parse
			case token.PING:
				c.Send(Pong, parse.Message)
			case token.PRIVMSG:
				c.handleEmotes(parse)

				displayName, ok := parse.Tags["display-name"]

				if !ok || displayName == "" {
					displayName = parse.Username
				}

				badges, err := c.handleBadges(parse)

				if err != nil {
					log.Printf("failed to get twitch badges with %s", err)
				}

				userID := parse.Tags["user-id"]

				c.RLock()
				cachedUser, ok := c.currentUsers[parse.Username]
				c.RUnlock()

				if !ok && userID != "" {
					twitchUser, err := c.twitchClient.GetUser(userID)
					if err != nil {
						log.Printf("failed to get user information with %s\n", err)
					} else {
						cachedUser = &user{
							profileImage: twitchUser.Logo,
						}
						c.Lock()
						c.currentUsers[parse.Username] = cachedUser
						c.Unlock()
					}
				}

				profileImage := ""

				if cachedUser != nil {
					profileImage = cachedUser.profileImage
				}

				msg := &Message{
					Message:      parse.Message,
					DisplayName:  displayName,
					Badges:       badges,
					ProfileImage: profileImage,
					Channel:      parse.Channel,
				}
				c.OnMessage <- msg
			}
		}
	}()
	return nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return errors.New("client is not started yet")
	}
	return c.conn.Close()
}

func (c *Client) Auth() error {
	err := c.capabilities()
	return err
}

func (c *Client) capabilities() error {
	err := c.Send(Cap, "twitch.tv/membership")

	if err != nil {
		return err
	}
	err = c.Send(Cap, "twitch.tv/tags")

	if err != nil {
		return err
	}

	err = c.Send(Cap, "twitch.tv/commands")

	if err != nil {
		return err
	}

	err = c.Send(Pass, c.conf.Auth)

	if err != nil {
		return err
	}

	err = c.Send(Nick, c.conf.Name)

	if err != nil {
		return err
	}

	err = c.Send(Join, c.conf.Channel)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Send(command chatCommand, message string) error {
	commandString, ok := commandToString[command]

	if !ok {
		return fmt.Errorf("unknown command %v", command)
	}

	if c.conn == nil {
		return errors.New("client is not started yet")
	}

	fmt.Fprintf(c.conn, "%s%s\r\n", commandString, message)
	return nil
}

func New(conf *Config) (*Client, error) {

	if conf == nil {
		return nil, fmt.Errorf("conf can't be nil")
	}

	if err := conf.validate(); err != nil {
		return nil, err
	}

	return &Client{
		conf:         conf,
		OnCap:        make(chan *parser.Message, 100),
		OnMessage:    make(chan *Message, 100),
		twitchEmotes: conf.TwitchEmotes,
		twitchClient: conf.TwitchAPI,
		badges:       conf.Badges,
		currentUsers: make(map[string]*user),
		emotesCache:  make(map[string]*emote),
	}, nil
}
