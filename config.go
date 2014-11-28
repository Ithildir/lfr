package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type config struct {
	LastUpdate int64
	Version    string
}

const configFileName string = "config.json"

func readConfig(homeDir string) (config, error) {
	path := filepath.Join(homeDir, configFileName)

	if !pathExists(path) {
		return config{}, nil
	}

	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return config{}, err
	}

	var cfg config

	err = json.Unmarshal(bytes, &cfg)

	return cfg, err
}

func (cfg config) save(homeDir string) error {
	path := filepath.Join(homeDir, configFileName)

	bytes, err := json.Marshal(cfg)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0777)
}
