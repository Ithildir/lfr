package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type config struct {
	Check int64
	Ver   int
}

const configFileName string = "config.json"

func readConfig(homeDir string) (config, error) {
	p := filepath.Join(homeDir, configFileName)

	if !fileExists(p) {
		return config{}, nil
	}

	b, err := ioutil.ReadFile(p)

	if err != nil {
		return config{}, err
	}

	var cfg config

	err = json.Unmarshal(b, &cfg)

	return cfg, err
}

func (cfg config) save(homeDir string) error {
	p := filepath.Join(homeDir, configFileName)

	b, err := json.Marshal(cfg)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(p, b, 0777)
}
