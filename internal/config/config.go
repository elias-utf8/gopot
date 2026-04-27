package config

import (
	"github.com/BurntSushi/toml"	
	"log"
)

type Config struct {
	Server ServerConfig `toml:"server"`
	Shell  ShellConfig  `toml:"shell"`
	Auth   AuthConfig   `toml:"auth"`
}

type ServerConfig struct {
	Port    int    `toml:"port"`
	Banner  string `toml:"banner"`
	HostKey string `toml:"host_key"`
}

type ShellConfig struct {
	Banner string `toml:"banner"`
	Prompt string `toml:"prompt"`
}

type AuthConfig struct {
	// MinAttempts: number of attempts to refuse per IP before accepting.
	// 0 = accept the first try, 2 = accept the 3rd attempt onward.
	MinAttempts int `toml:"min_attempts"`
}

func LoadConfig() *Config {
	var cfg Config
	_, err := toml.DecodeFile("gopot.toml", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
