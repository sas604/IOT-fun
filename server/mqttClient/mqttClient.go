package mqttclient

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sas604/IOT-fun/server/db"
)

var Client mqtt.Client

type SensorData struct {
	Temp float64 `json:"temp"`
	Hum  float64 `json:"hum"`
	CO   float64 `json:"co"`
}

var msgH mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	var mes SensorData
	err := json.Unmarshal([]byte(m.Payload()), &mes)
	if err != nil {
		fmt.Print(err.Error())

	}
	writeApi := db.DB.WriteAPI("me", "iot-fun")
	writeApi.WriteRecord(fmt.Sprintf("sht-31,controller=main,sensor=sht-31,location=%s temp=%f,hum=%f,co=%f", "any", mes.Temp, mes.Hum, mes.CO))
	writeApi.Flush()
}

func ConectToMqttBroker() {
	options := mqtt.NewClientOptions().AddBroker("tcp://192.168.1.106:1883")
	options.SetUsername("pi")
	options.SetPassword("boopyou")
	options.SetOrderMatters(false)
	options.OnConnect = func(c mqtt.Client) {
		c.Subscribe("mush/sensor-group/mesurments", 0, msgH)
	}
	Client = mqtt.NewClient(options)

	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to mqtt broker \n")
	}
}
