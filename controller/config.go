package controller

import (
	"os"

	"gopkg.in/yaml.v3"
)

type TV_CONFIG struct {
	TurnServerAddr string `yaml:"turnServerAddr"`
}

var TVCONFIG *TV_CONFIG = nil

func LoadConfig() error {

	tv_config := TV_CONFIG{}

	file_b, err := os.ReadFile("./config.yaml")

	if err != nil {

		return err
	}

	err = yaml.Unmarshal(file_b, &tv_config)

	if err != nil {

		return err
	}

	TVCONFIG = &tv_config

	return nil

}
