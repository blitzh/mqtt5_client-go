package main

import "time"

var (
	mqServer         string        = "mqtt://broker.emqx.io:1883"
	mqClientID       string        = "MQTT5CLient" // Change this to something random if using a public test server
	mqTopicSubs      []string      = []string{"MQTT5CLient/sub"}
	mqTopicPub       string        = "MQTT5CLient/pub"
	mqUsername       string        = "MasGanteng"
	mqPassword       string        = "Blitar"
	mqKeepAlive      uint16        = 60
	mqQOS            byte          = 0
	mqRetained       bool          = false
	mqCleanStart     bool          = true
	mqConnRetryDelay time.Duration = 30 * time.Second
	mqConnTimeout    time.Duration = 30 * time.Second
	mqSessionExpiry  uint32        = 7200 //7200 // 2 hours

	//Status
	mqConnected bool = false
)
