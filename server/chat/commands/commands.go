package commands

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/miguel250/streaming-setup/server/irc"
)

type Command struct {
	Description string `json:"description"`
	AllowRoles  AllowRoles
	Action      CommandFunc
}

type AllowRoles []string

func (a AllowRoles) Allow(msg *irc.Message) error {
	hasPermission := false

	for _, badge := range msg.Badges {
		b := strings.ToLower(badge.Title)
		if a.Has(b) {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return fmt.Errorf("user is not allow to use command: %s", msg.DisplayName)
	}
	return nil
}

func (a AllowRoles) Has(role string) bool {
	for _, allow := range a {
		if allow == role {
			return true
		}
	}
	return false
}

type CommandFunc func(client *irc.Client, msg *irc.Message, allowRoles AllowRoles) error

type AvailableCommands struct {
	sync.RWMutex
	conf     *Config
	client   *irc.Client
	commands map[string]*Command
	shutdown chan struct{}
}

func (a *AvailableCommands) Start() {
	log.Println("Handling chat commands")
	messageChannel := a.client.MessageListener()
	go func() {
		for {
			select {
			case msg := <-messageChannel:
				err := a.parseMsg(msg)
				if err != nil {
					log.Println(err)
				}
			case <-a.shutdown:
				return
			}
		}
	}()
}

func (a *AvailableCommands) parseMsg(msg *irc.Message) error {
	if msg.Message[0] != '!' {
		return nil
	}

	msgSlice := strings.Split(msg.Message, " ")
	command := msgSlice[0][1:]

	a.RLock()
	val, ok := a.commands[command]
	a.RUnlock()
	if !ok {
		return fmt.Errorf("invalid command %s", command)
	}

	return val.Action(a.client, msg, val.AllowRoles)
}

func (a *AvailableCommands) Close() {
	a.shutdown <- struct{}{}
}

func (a *AvailableCommands) printHelpCommand() *Command {
	return &Command{
		Description: "Print all chat bot commands",
		Action: func(client *irc.Client, msg *irc.Message, allowRoles AllowRoles) error {
			hiMsg := "Hi, here is a list of commands"
			err := client.SendMessage(hiMsg)
			if err != nil {
				log.Printf("Unable to send message with %s\n", err)
			}

			for key, value := range a.commands {
				msg := fmt.Sprintf("- !%s - %s", key, value.Description)
				err := client.SendMessage(msg)
				if err != nil {
					log.Printf("Failed to send help command with %s\n", err)
				}
			}
			return nil
		},
	}
}

func (a *AvailableCommands) shoutout() *Command {
	return &Command{
		Description: "Give a shoutout to someone",
		AllowRoles: []string{
			"broadcaster",
			"moderator",
		},
		Action: func(client *irc.Client, msg *irc.Message, allowRoles AllowRoles) error {

			if err := allowRoles.Allow(msg); err != nil {
				return fmt.Errorf("user is not allow to use command: %s", msg.DisplayName)
			}

			msgSlice := strings.Split(msg.Message, " ")

			if len(msgSlice) == 1 {
				client.SendMessage("Missing username @example")
				return nil
			}

			username := msgSlice[1]
			if username[0] != '@' {
				client.SendMessage("username doesn't include @")
				return nil
			}

			shoutoutMsg := fmt.Sprintf("Go checkout - http://twitch.tv/%s", username[1:])
			client.SendMessage(shoutoutMsg)
			return nil
		},
	}
}

// !addcmd discord - description - Please join our discord server - https://discord.gg/3q2vkv
func (a *AvailableCommands) addcmd() *Command {
	return &Command{
		Description: "Add a new command to chat bot",
		AllowRoles: []string{
			"broadcaster",
			"moderator",
		},
		Action: func(client *irc.Client, msg *irc.Message, allowRoles AllowRoles) error {
			if err := allowRoles.Allow(msg); err != nil {
				return fmt.Errorf("user is not allow to use command: %s", msg.DisplayName)
			}

			msgSlice := strings.Split(msg.Message, "-")

			if len(msgSlice) < 3 {
				client.SendMessage("addcmd needs 3 args marked by '-'")
				client.SendMessage("- !addcmd discord - description - Please join our discord server - url")
				return nil
			}

			commandSlice := strings.Split(strings.Trim(msgSlice[0], " "), " ")
			if len(commandSlice) != 2 {
				client.SendMessage("Unable to parse command name")
				return nil
			}

			commandName := commandSlice[1]
			description := strings.Trim(msgSlice[1], " ")
			message := strings.Trim(strings.Join(msgSlice[2:], "-"), " ")

			a.AddCommand(commandName, message, description)
			a.conf.AddCommand(commandName, CommandConfig{
				Description: description,
				Message:     message,
			})

			err := a.conf.Save()
			if err != nil {
				client.SendMessage("Failed to save command")
				return fmt.Errorf("failed to save command with %s", err)
			}

			client.SendMessage(fmt.Sprintf("Command (!%s - %s - %s) was added successfully.", commandName, description, message))
			return nil
		},
	}
}

func (a *AvailableCommands) AddCommand(cmd, message, description string) *Command {
	action := func(client *irc.Client, _ *irc.Message, _ AllowRoles) error {
		err := client.SendMessage(message)
		if err != nil {
			log.Printf("Failed to %s help command with %s\n", cmd, err)
		}
		return nil
	}

	a.Lock()
	defer a.Unlock()
	command := &Command{
		Description: description,
		Action:      action,
	}
	a.commands[cmd] = command
	return command
}

func New(client *irc.Client, conf *Config) *AvailableCommands {
	available := &AvailableCommands{
		conf:     conf,
		client:   client,
		commands: make(map[string]*Command),
		shutdown: make(chan struct{}),
	}

	available.commands["commands"] = available.printHelpCommand()
	available.commands["so"] = available.shoutout()
	available.commands["addcmd"] = available.addcmd()

	for key, value := range conf.Commands {
		available.AddCommand(key, value.Message, value.Description)
	}
	return available
}
