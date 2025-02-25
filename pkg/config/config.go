package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Rummage struct {
		Version      string `toml:"version"`
		DbApiVersion string `toml:"dbApiVersion"`
	} `toml:"rummage"`
}

func SetVersions() *Config {
	var conf Config
	if _, err := toml.DecodeFile("./pkg/config/version.toml", &conf); err != nil {
		panic(err)
	}

	return &conf
}
