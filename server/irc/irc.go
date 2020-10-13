package irc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"strconv"
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
	connMutex      sync.RWMutex
	conf           *Config
	conn           net.Conn
	OnCap          chan *parser.Message
	onMessages     []chan *Message
	onClearMessage []chan *ClearMessage
	OnReconnect    chan bool
	twitchEmotes   *twitchemotes.API
	twitchClient   *twitch.API
	badges         map[string]*twitch.BadgeVersion
	currentUsers   map[string]*user
	emotesCache    map[string]*emote
	reader         *textproto.Reader
}

type Message struct {
	Badges       []*twitch.Badge `json:"badges"`
	DisplayName  string          `json:"display-name"`
	Message      string          `json:"message"`
	ProfileImage string          `json:"profile_image"`
	Channel      string          `json:"channel"`
}

type ClearMessage struct {
	Message   string
	UserLogin string
	Channel   string
	MessageID string
	Timestamp int64
}

type user struct {
	profileImage string
}

type emote struct {
	id       string
	code     string
	imageURL string
}

func (c *Client) Start() error {
	err := c.connect()
	if err != nil {
		return err
	}

	go func() {
		for {
			line, err := c.reader.ReadLine()

			if err != nil && err == io.EOF {
				log.Println("Connection was closed")
				return
			}

			if err != nil {
				return
			}

			parse, err := parser.ParseMsg(line)

			if err != nil {
				log.Printf("Failed to parse message with %s\n", err)
			}

			switch parse.Command {
			case token.CAP:
				c.OnCap <- parse
			case token.RECONNECT:
				err := c.connect()
				if err != nil {
					log.Printf("failed to reconnect to server with: %s", err)
					return
				}
				c.OnReconnect <- true
			case token.PING:
				err := c.Send(Pong, parse.Message)
				if err != nil {
					log.Printf("failed to send pong command to server")
				}
			case token.CLEARMSG:
				msg := &ClearMessage{
					Message:   parse.Message,
					UserLogin: parse.Tags["login"],
					Channel:   parse.Channel,
					MessageID: parse.Tags["target-msg-id"],
				}

				if i, err := strconv.ParseInt(parse.Tags["tmi-sent-ts"], 10, 64); err == nil {
					msg.Timestamp = i
				}

				c.RLock()
				for _, channel := range c.onClearMessage {
					select {
					case channel <- msg:
					default:
					}
				}
				c.RUnlock()
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

				c.RLock()
				for _, channel := range c.onMessages {
					channel <- msg
				}
				c.RUnlock()
			}
		}
	}()
	return nil
}

func (c *Client) connect() error {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	conn, err := net.Dial("tcp", c.conf.URL)
	if err != nil {
		return fmt.Errorf("failed to create connection with %w", err)
	}

	oldConn := c.conn
	c.conn = conn
	if oldConn != nil {
		oldConn.Close()
	}

	reader := bufio.NewReader(c.conn)
	c.reader = textproto.NewReader(reader)

	return nil
}

func (c *Client) Close() error {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	if c.conn == nil {
		return errors.New("client is not started yet")
	}

	c.Lock()
	defer c.Unlock()

	// close all on message channels
	for _, channel := range c.onMessages {
		close(channel)
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

func (c *Client) SendMessage(msg string) error {
	return c.Send(PrivMsg, fmt.Sprintf("%s :%s", c.conf.Channel, msg))
}

func (c *Client) MessageListener() chan *Message {
	channel := make(chan *Message)
	c.Lock()
	defer c.Unlock()
	c.onMessages = append(c.onMessages, channel)
	return channel
}

func (c *Client) ClearMessageListener() chan *ClearMessage {
	channel := make(chan *ClearMessage)
	c.Lock()
	defer c.Unlock()
	c.onClearMessage = append(c.onClearMessage, channel)
	return channel
}

func (c *Client) Send(command chatCommand, message string) error {
	commandString, ok := commandToString[command]

	if !ok {
		return fmt.Errorf("unknown command %v", command)
	}

	c.connMutex.Lock()
	defer c.connMutex.Unlock()

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
		conf:           conf,
		OnCap:          make(chan *parser.Message, 100),
		onMessages:     make([]chan *Message, 0, 10),
		onClearMessage: make([]chan *ClearMessage, 0, 10),
		OnReconnect:    make(chan bool, 10),
		twitchEmotes:   conf.TwitchEmotes,
		twitchClient:   conf.TwitchAPI,
		badges:         conf.Badges,
		currentUsers:   make(map[string]*user),
		emotesCache:    make(map[string]*emote),
	}, nil
}
