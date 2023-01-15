package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var msgH MQTT.MessageHandler = func(c MQTT.Client, m MQTT.Message) {
	if m.Topic() == "test/temperature" {
		fmt.Printf("Temperature: %s C\n", m.Payload())
	} else {
		fmt.Printf("Humidity: %s %%\n", m.Payload())
	}

}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	options := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	options.SetUsername("pi")
	options.SetPassword("boopyou")
	options.SetOrderMatters(false)
	options.OnConnect = func(c MQTT.Client) {
		c.Subscribe("test/temperature", 0, msgH)
		c.Subscribe("test/humidity", 0, msgH)
	}
	client := MQTT.NewClient(options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server \n")
	}
	<-c

}
