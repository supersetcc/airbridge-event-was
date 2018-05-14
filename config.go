package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	PRODUCTION_PATH = "res/config/production.yml"
	DEVELOP_PATH    = "res/config/develop.yml"
)

type Config struct {
	Kafka struct {
		BrokerList   []string `yaml:"kafka_broker_list"`
		ProduceTopic string   `yaml:"kafka_produce_topic"`
	}
}

func LoadConfig() (*Config, error) {
	stream, err := ioutil.ReadFile(PRODUCTION_PATH)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := yaml.Unmarshal(stream, config); err != nil {
		return nil, err
	}

	return &config, nil
}
