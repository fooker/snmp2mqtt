package main

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"log"
	"strings"
)

type MQTT struct {
	c     mqtt.Client
	realm string
}

type Status string

func MQTTConnect(config *MQTTConfig) *MQTT {
	if !strings.Contains(config.Address, ":") {
		config.Address += ":1883"
	}

	realm := config.Realm
	if realm != "" && !strings.HasSuffix(realm, "/") {
		realm += "/"
	}

	c := mqtt.NewClient(mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s/", config.Address)).
		SetClientID(config.ClientID))

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT: Failed to connect to broker: %v", token.Error())

	} else {
		log.Printf("MQTT: Connected to broker: %s", config.Address)
	}

	return &MQTT{
		c:     c,
		realm: realm,
	}
}

func (m *MQTT) Close() {
	m.c.Disconnect(0)
}

func (m *MQTT) Publish(topic string, message string) {
	topic = m.realm + topic

	if token := m.c.Publish(topic, 0, true, []byte(message)); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT: Failed to publish: %s: %v", topic, token.Error())
	}
}
