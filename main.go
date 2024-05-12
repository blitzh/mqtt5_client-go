package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

var mqGlobalCfg autopaho.ClientConfig
var mqClient5 *autopaho.ConnectionManager

func SetInitialConfig(topicForSubs []string) {

	serverUrl, err := url.Parse(mqServer)
	if err != nil {
		panic(err)
	}
	mqGlobalCfg = autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{serverUrl},
		ConnectUsername:               mqUsername,
		ConnectPassword:               []byte(mqPassword),
		KeepAlive:                     mqKeepAlive,
		ConnectRetryDelay:             mqConnRetryDelay,
		ConnectTimeout:                mqConnTimeout,
		CleanStartOnInitialConnection: true,
		SessionExpiryInterval:         mqSessionExpiry,

		WillMessage: &paho.WillMessage{
			Topic:   mqTopicPub,
			QoS:     mqQOS,
			Payload: []byte("Offline#" + time.Now().UTC().String()),
			Retain:  mqRetained,
		},

		// WillProperties: &paho.WillProperties{
		// 	WillDelayInterval: 0,
		// 	PayloadFormatIndicator: paho.PayloadFormatIndicatorUTF8,

		// },
		// WillProperties *paho.WillProperties
		// TlsCfg                        *tls.Config
		// WebSocketCfg      	  *websocket.Config

		//AttemptConnection bool

		OnConnectError: func(err error) {
			mqConnected = false
			fmt.Printf("error whilst attempting connection: %s\n", err)
		},

		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			mqConnected = true
			fmt.Println("mqtt connection up")
			// Subscribing in the OnConnectionUp callback is recommended (ensures the subscription is reestablished if
			// the connection drops)

			listTopic := []paho.SubscribeOptions{}
			for _, topic := range topicForSubs {
				listTopic = append(listTopic, paho.SubscribeOptions{Topic: topic, QoS: mqQOS})
			}

			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: listTopic,
			}); err != nil {
				fmt.Printf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
			}
			fmt.Println("mqtt subscription made")
		},
	}
}

func PublishMsg(ctx context.Context, topic string, msg string) {
	if mqConnected {
		resp, err := mqClient5.Publish(ctx, &paho.Publish{
			QoS:     mqQOS,
			Topic:   topic,
			Payload: []byte("response " + msg),
		})
		if err != nil {
			fmt.Println(err) // Publish will exit when context cancelled or if something went wrong
		}
		fmt.Println("Response: ", &resp)
	}
}

func startMQTT5() {
	ctx, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelSignal()

	SetInitialConfig(mqTopicSubs)

	router := paho.NewStandardRouter()
	router.DefaultHandler(func(p *paho.Publish) { fmt.Printf("defaulthandler received message with topic: %s\n", p.Topic) })

	mqGlobalCfg.ClientConfig = paho.ClientConfig{
		ClientID: mqClientID,
		OnPublishReceived: []func(paho.PublishReceived) (bool, error){
			func(pr paho.PublishReceived) (bool, error) {
				// fmt.Println("<<- Recv | topic: ", pr.Packet.Topic)
				// fmt.Println("Payload: ", string(pr.Packet.Payload))
				router.Route(pr.Packet.Packet())
				return true, nil // we assume that the router handles all messages (todo: amend router API)
			}},

		OnClientError: func(err error) { fmt.Printf("client error: %s\n", err) },
		OnServerDisconnect: func(d *paho.Disconnect) {
			mqConnected = false
			if d.Properties != nil {
				fmt.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
			} else {
				fmt.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
			}
		},
	}

	var err error
	mqClient5, err = autopaho.NewConnection(ctx, mqGlobalCfg)
	if err != nil {
		fmt.Println(err)
	}

	if err = mqClient5.AwaitConnection(ctx); err != nil {
		fmt.Println(err)
	}

	router.RegisterHandler(mqTopicSubs[0], func(p *paho.Publish) {
		fmt.Printf("topic: %s\n", p.Topic)
		fmt.Println("Payload: ", string(p.Payload))
		//Send back utk tes kepastian
		PublishMsg(ctx, mqTopicPub, "Received: "+string(p.Payload))
	})

	// Publish Online First
	PublishMsg(ctx, mqTopicPub, "Online#"+time.Now().UTC().String())

	// ticker := time.NewTicker(time.Duration(mqKeepAlive) * time.Second)
	// defer ticker.Stop()
	// for {
	// 	select {
	// 	case <-ticker.C:
	// 		//CHeck if the connection is still alive
	// 		if mqConnected {
	// 			fmt.Println("Connected")
	// 		} else {
	// 			fmt.Println("Disconnected")
	// 		}
	// 		continue
	// 	case <-ctx.Done():
	// 	}
	// 	break
	// }
	select {
	//case error
	case <-ctx.Done():
		fmt.Println("Context Done")
		//case disconet
	case <-mqClient5.Done():
		fmt.Println("MQTT5 Done")

	}
}

func InitMQTT5() {
	startMQTT5()
}

func main() {

	go InitMQTT5()

	// make alwasy on running
	select {}
}
