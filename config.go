package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	PRODUCTION_PATH = "res/config/production.yml"
	DEVELOP_PATH    = "res/config/develop.yml"
)

type Config struct {
	Kafka struct {
		BrokerList []string `yaml:"broker_list"`
	} `yaml:"kafka"`

	Newrelic struct {
		License string `yaml:"license"`
		AppName string `yaml:"app_name"`
	} `yaml:"newrelic"`

	Server struct {
		Port int `yaml:"server"`
	} `yaml:"port"`
}

func LoadConfig(production bool) (*Config, error) {
	configPath := DEVELOP_PATH
	if production {
		configPath = PRODUCTION_PATH
	}

	stream, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := yaml.Unmarshal(stream, &config); err != nil {
		return nil, err
	}

	if len(config.Kafka.BrokerList) == 0 {
		return nil, fmt.Errorf("kafka brokers must over than one")
	}

	return &config, nil
}
