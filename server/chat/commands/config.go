package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type Config struct {
	mux      sync.Mutex               `json:"-"`
	path     string                   `json:"-"`
	Commands map[string]CommandConfig `json:"commands"`
}

type CommandConfig struct {
	Description string `json:"description"`
	Message     string `json:"message"`
}

func (c *Config) AddCommand(name string, cmd CommandConfig) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.Commands == nil {
		c.Commands = make(map[string]CommandConfig)
	}
	c.Commands[name] = cmd
}

func (c *Config) Save() error {
	c.mux.Lock()
	defer c.mux.Unlock()
	b, err := json.MarshalIndent(c, "", "   ")
	if err != nil {
		return fmt.Errorf("failed to save configuration with %s", err)
	}

	err = ioutil.WriteFile(c.path, b, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file with %s", err)
	}

	return nil
}

func NewConfig(path string) (*Config, error) {
	body, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("failed to get configuration file with %w", err)
	}

	conf := &Config{path: path}
	err = json.Unmarshal(body, conf)

	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration file with %w", err)
	}

	return conf, nil
}
