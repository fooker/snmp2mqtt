package main

import (
	"flag"
	"log"
	"time"
	snmp "github.com/soniah/gosnmp"
	"fmt"
)

var intervalFlag = flag.Duration("interval", 5*time.Minute, "Scraping interval")
var configFlag = flag.String("config", "config.json", "Path to config file")

type Host struct {
	name string

	agent *snmp.GoSNMP

	querynames []string
	queryoids  []string
}

func main() {
	// Parse and check flags
	flag.Parse()

	// Load and parse config
	config, err := LoadConfig(*configFlag)
	if err != nil {
		log.Fatalf("Failed to load config: %s: %v", *configFlag, err)
	}

	// Create host objects
	hosts := make([]*Host, 0)
	for host, config := range config.Hosts {
		if config.Port == 0 {
			config.Port = 161
		}

		var agent *snmp.GoSNMP
		switch config.Version {
		case "1":
			agent = &snmp.GoSNMP{
				Target:    config.Host,
				Port:      config.Port,
				Community: config.Community,
				Version:   snmp.Version1,
				Timeout:   time.Duration(2) * time.Second,
			}

		case "2c":
			agent = &snmp.GoSNMP{
				Target:    config.Host,
				Port:      config.Port,
				Community: config.Community,
				Version:   snmp.Version2c,
				Timeout:   time.Duration(2) * time.Second,
			}
		}

		querynames := make([]string, 0)
		queryoids := make([]string, 0)
		for name, oid := range config.OIDs {
			querynames = append(querynames, name)
			queryoids = append(queryoids, oid)
		}

		hosts = append(hosts, &Host{
			name:       host,
			agent:      agent,
			querynames: querynames,
			queryoids:  queryoids,
		})
	}

	// Connect to MQTT
	mqtt := MQTTConnect(config.MQTT)
	defer mqtt.Close()

	log.Printf("Connected...")

	// Interval loop
	for range time.Tick(*intervalFlag) {
		for _, host := range hosts {
			err = host.agent.Connect()
			if err != nil {
				log.Fatalf("SNMP: Connect to %s failed: %v", host.name, err)
			}

			result, err := host.agent.Get(host.queryoids)
			if err != nil {
				log.Fatalf("SNMP: Error while getting values: %v", err)
			}

			for i, variable := range result.Variables {
				log.Printf("%s/%s = %v", host.name, host.querynames[i], variable.Value)
				mqtt.Publish(fmt.Sprintf("%s/%s", host.name, host.querynames[i]), fmt.Sprintf("%v", variable.Value))
			}

			host.agent.Conn.Close()
		}
	}
}
