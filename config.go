package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type tomlConfig struct {
	Bucket      string          `toml:"bucket"`
	Region      string          `toml:"region"`
	Prefix      string          `toml:"prefix"`
	Mode        string          `toml:"mode"`
	FilePattern string          `toml:"filepattern"`
	Rules       map[string]rule //`toml:"rules"`
}

type rule struct {
	Rule string `toml:"rule"`
}

var verbose = false
var sandbox = false

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
