package commands

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/miguel250/streaming-setup/server/irc"
	"github.com/miguel250/streaming-setup/server/irc/util"
)

var update = flag.Bool("update", false, "update golden files")

func TestCommands(t *testing.T) {
	for _, test := range []struct {
		name           string
		messageCount   int
		inputFileName  string
		outputFileName string
		errorMessage   string
		config         *Config
	}{
		{
			"help command",
			4,
			"help_command_message.json",
			"help_command_result.json",
			"",
			&Config{},
		},
		{
			"so command no permissions",
			0,
			"so_command_message_missing_permissions.json",
			"",
			"user is not allow to use command: AttackKopter",
			&Config{},
		},
		{
			"so command missing username",
			1,
			"so_command_message_missing_user.json",
			"so_command_message_missing_result_user.json",
			"",
			&Config{},
		},
		{
			"so command",
			1,
			"so_command_message.json",
			"so_command_message_result.json",
			"",
			&Config{},
		},
		{
			"discord command",
			1,
			"discord_command_message.json",
			"discord_command_result.json",
			"",
			&Config{
				Commands: map[string]CommandConfig{
					"discord": {
						Description: "Print discord server URL",
						Message:     "Please join our discord server - https://discord.gg/3q2vkv",
					},
				},
			},
		},
		{
			"addcmd command",
			1,
			"addcmd_command_message.json",
			"addcmd_command_result.json",
			"",
			&Config{},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			client, _ := util.CreateMockChatClient(t)
			tmpfile, err := ioutil.TempFile("", fmt.Sprintf("commands_%s.json", test.name))
			if err != nil {
				t.Fatalf("failed to create tmp configuration file with %s", err)
			}

			defer os.Remove(tmpfile.Name())
			test.config.path = tmpfile.Name()

			b, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s", test.inputFileName))
			if err != nil {
				t.Fatalf("Failed to run test with %s", err)
			}

			msg := &irc.Message{}
			err = json.Unmarshal(b, msg)
			if err != nil {
				t.Fatalf("Failed to unmarshal json with %s", err)
			}

			client.Start()

			commands := New(client, test.config)
			commands.Start()
			defer commands.Close()

			err = commands.parseMsg(msg)
			if err != nil && test.errorMessage == "" {
				t.Fatalf("Failed to parse message with %s", err)
			}

			if test.errorMessage != "" {
				if err == nil {
					t.Fatalf("expected an error to be trigger")
				}

				if test.errorMessage != err.Error() {
					t.Fatalf("Expected an error but error message didn't match got: '%s', want: '%s'", err, test.errorMessage)
				}
				return
			}

			msgChannel := client.MessageListener()
			receiveMessages := make([]*irc.Message, 0, test.messageCount)

			for {
				msg := <-msgChannel
				receiveMessages = append(receiveMessages, msg)

				test.messageCount--
				if test.messageCount == 0 {
					break
				}
			}

			testCompareGoldenFiles(fmt.Sprintf("testdata/%s", test.outputFileName), receiveMessages, t)
		})
	}
}

func TestDynamicAddedCommand(t *testing.T) {
	client, _ := util.CreateMockChatClient(t)
	tmpfile, err := ioutil.TempFile("", "commands_dynamic.json")
	if err != nil {
		t.Fatalf("failed to create tmp configuration file with %s", err)
	}

	defer os.Remove(tmpfile.Name())
	conf := &Config{path: tmpfile.Name()}

	commands := New(client, conf)
	commands.Start()
	defer commands.Close()
	client.Start()

	b, err := ioutil.ReadFile("testdata/addcmd_command_message.json")
	if err != nil {
		t.Fatalf("Failed to run test with %s", err)
	}

	msg := &irc.Message{}
	err = json.Unmarshal(b, msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal json with %s", err)
	}

	err = commands.parseMsg(msg)
	if err != nil {
		t.Fatalf("Failed to parse message with %s", err)
	}

	msgChannel := client.MessageListener()
	fmt.Println(<-msgChannel)

	b, err = ioutil.ReadFile("testdata/discord_command_message.json")
	if err != nil {
		t.Fatalf("Failed to run test with %s", err)
	}

	msg = &irc.Message{}
	err = json.Unmarshal(b, msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal json with %s", err)
	}

	err = commands.parseMsg(msg)
	if err != nil {
		t.Fatalf("Failed to parse message with %s", err)
	}

	receiveMessages := <-msgChannel
	fmt.Println(receiveMessages.Message)
	testCompareGoldenFiles("testdata/dynamic_discord_command_result.json", receiveMessages, t)
}

func testCompareGoldenFiles(wantFilename string, v interface{}, t *testing.T) {
	got, err := json.Marshal(v)

	if err != nil {
		t.Fatalf("Failed to marshal to json with %s", err)
	}

	if *update {
		err := ioutil.WriteFile(wantFilename, got, 0644)
		if err != nil {
			t.Fatalf("Failed to save golden file with %s", err)
		}
	}

	want, err := ioutil.ReadFile(wantFilename)

	if err != nil {
		t.Fatalf("Failed to open golden file with %s", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf("Doesn't match got (%s), want (%s)", string(got), string(want))
	}
}
