package main

// import (
// 	"context"
// 	"fmt"
// 	"net/url"
// 	"os"
// 	"os/signal"
// 	"strconv"
// 	"syscall"
// 	"time"

// 	"github.com/eclipse/paho.golang/autopaho"
// 	"github.com/eclipse/paho.golang/paho"
// )

// func main() {
// 	// App will run until cancelled by user (e.g. ctrl-c)
// 	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
// 	defer stop()

// 	// We will connect to the Eclipse test server (note that you may see messages that other users publish)
// 	u, err := url.Parse("mqtt://mqx1.jtisrv.com:1883")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	cliCfg := autopaho.ClientConfig{
// 		ServerUrls:                    []*url.URL{u},
// 		ConnectUsername:               mqUsername,
// 		ConnectPassword:               []byte(mqPassword),
// 		KeepAlive:                     mqKeepAlive,
// 		CleanStartOnInitialConnection: false,
// 		// (60 = 1 minute, 3600 = 1 hour, 86400 = one day, 0xFFFFFFFE = 136 years, 0xFFFFFFFF = don't expire)
// 		SessionExpiryInterval: 7200,
// 		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
// 			fmt.Println("mqtt connection up")
// 			// Subscribing in the OnConnectionUp callback is recommended (ensures the subscription is reestablished if
// 			// the connection drops)
// 			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
// 				Subscriptions: []paho.SubscribeOptions{
// 					{Topic: topic, QoS: mqQOS},
// 				},
// 			}); err != nil {
// 				fmt.Printf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
// 			}
// 			fmt.Println("mqtt subscription made")
// 		},
// 		OnConnectError: func(err error) { fmt.Printf("error whilst attempting connection: %s\n", err) },
// 		// eclipse/paho.golang/paho provides base mqtt functionality, the below config will be passed in for each connection
// 		ClientConfig: paho.ClientConfig{
// 			// If you are using QOS 1/2, then it's important to specify a client id (which must be unique)
// 			ClientID: clientID,
// 			// OnPublishReceived is a slice of functions that will be called when a message is received.
// 			// You can write the function(s) yourself or use the supplied Router
// 			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
// 				func(pr paho.PublishReceived) (bool, error) {
// 					fmt.Printf("received message on topic %s; body: %s (retain: %t)\n", pr.Packet.Topic, pr.Packet.Payload, pr.Packet.Retain)
// 					return true, nil
// 				}},
// 			OnClientError: func(err error) { fmt.Printf("client error: %s\n", err) },
// 			OnServerDisconnect: func(d *paho.Disconnect) {
// 				if d.Properties != nil {
// 					fmt.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
// 				} else {
// 					fmt.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
// 				}
// 			},
// 		},
// 	}

// 	c, err := autopaho.NewConnection(ctx, cliCfg) // starts process; will reconnect until context cancelled
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// Wait for the connection to come up
// 	if err = c.AwaitConnection(ctx); err != nil {
// 		fmt.Println(err)
// 	}

// 	ticker := time.NewTicker(time.Second)
// 	msgCount := 0
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
// 			msgCount++
// 			// Publish a test message (use PublishViaQueue if you don't want to wait for a response)
// 			if _, err = c.Publish(ctx, &paho.Publish{
// 				QoS:     1,
// 				Topic:   topic,
// 				Payload: []byte("TestMessage: " + strconv.Itoa(msgCount)),
// 			}); err != nil {
// 				if ctx.Err() == nil {
// 					fmt.Println(err) // Publish will exit when context cancelled or if something went wrong
// 				}
// 			}
// 			continue
// 		case <-ctx.Done():
// 		}
// 		break
// 	}

// 	fmt.Println("signal caught - exiting")
// 	<-c.Done() // Wait for clean shutdown (cancelling the context triggered the shutdown)
// }
