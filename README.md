snmp2mqtt
=========
`snmp2mqtt` is a daemon polling a set of SNMP OIDs and export the resulting values to MQTT.

Usage
-----
The daemon must be started with a config file containing all the thing to query.
At first one must specify the hosts to query including the credentials and settings used to communicate with this host via SNMP.
Second, for each host one can specify a mapping from exported topic names to OIDs to query from that specific host.
 
Adopt the `config.yaml.example` to your needs.

Limitations
-----------
For now, this is very simple and just supports plain OIDs. Neither tables nor sets are supported.
It is also limited to SNMP version 1 and 2c for now.
