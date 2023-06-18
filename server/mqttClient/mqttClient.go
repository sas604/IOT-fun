package mqttclient

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var Client mqtt.Client

func ConectToMqttBroker() {
	options := mqtt.NewClientOptions().AddBroker("tcp://192.168.1.106:1883")
	options.SetUsername("pi")
	options.SetPassword("boopyou")
	options.SetOrderMatters(false)
	// options.OnConnect = func(c MQTT.Client) {
	// 	c.Subscribe("mush/sensor-group/mesurments", 0, msgH)
	// }
	Client = mqtt.NewClient(options)

	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to mqtt broker \n")
	}
}
