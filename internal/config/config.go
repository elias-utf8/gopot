package config

import (
	"github.com/BurntSushi/toml"	
	"log"
)

type Config struct {
	Server ServerConfig `toml:"server"`
	Shell ShellConfig `toml:"shell"`
}

type ServerConfig struct {
	Port int `toml:"port"`
	Banner string `toml:"banner"`
}

type ShellConfig struct {
	Banner string `toml:"banner"`
	Prompt string `toml:"prompt"`
}

func LoadConfig() *Config {
	var cfg Config
	_, err := toml.DecodeFile("gopot.toml", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
