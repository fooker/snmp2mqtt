package main

import (
	"gopkg.in/yaml.v2"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type MQTTConfig struct {
	Address  string `yaml:"address"`
	ClientID string `yaml:"client_id"`
	Realm    string `yaml:"realm"`
}

type HostConfig struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`

	Version string `yaml:"version"`
	Community string `yaml:"community"`

	OIDs map[string]string `yaml:"oids"`
}

type Config struct {
	Interval time.Duration `yaml:"interval"`

	MQTT *MQTTConfig `yaml:"mqtt"`

	Hosts map[string]HostConfig `yaml:"hosts"`
}

func LoadConfig(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
	}

	config := &Config{}
	if err := yaml.Unmarshal(b, config); err != nil {
		return nil, err
	}

	log.Printf("Config: Loaded: %v", config)

	return config, nil
}
