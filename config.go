package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type tomlConfig struct {
	Bucket string          `toml:"bucket"`
	Region string          `toml:"region"`
	Prefix string          `toml:"prefix"`
	Rules  map[string]rule //`toml:"rules"`
}

type rule struct {
	Rule string `toml:"rule"`
}

func readConfig(fname string) (tomlConfig, error) {
	var cfg tomlConfig
	content, err := ioutil.ReadFile(fname)

	if err != nil {
		return cfg, err
	}

	if _, err := toml.Decode(string(content), &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
