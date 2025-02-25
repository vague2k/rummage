package config

import (
	"github.com/BurntSushi/toml"
)

type config struct {
	Rummage struct {
		Version      string `toml:"version"`
		DbApiVersion string `toml:"dbApiVersion"`
	} `toml:"rummage"`
}

func SetVersions() *config {
	var conf config
	if _, err := toml.DecodeFile("./pkg/config/version.toml", &conf); err != nil {
		panic(err)
	}

	return &conf
}
