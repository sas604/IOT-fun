package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SensorData struct {
	Temp     float64 `json:"temp"`
	Hum      float64 `json:"hum"`
	CO       float64 `json:"co"`
	TubeHum  float64 `json:"tube_hum"`
	TubeTemp float64 `json:"tube_temp"`
}

func (app *application) hadleIncomingMeasurements(c mqtt.Client, m mqtt.Message) {
	// var mes SensorData
	// err := json.Unmarshal([]byte(m.Payload()), &mes)
	// if err != nil {
	// 	fmt.Print(err.Error())

	// }
	// writeApi := db.DB.WriteAPI("me", "iot-fun")
	// writeApi.WriteRecord(fmt.Sprintf("sht-31,sensor=sht-31 temp=%f,hum=%f,co=%f,tube_hum=%f,tube_temp=%f", mes.Temp, mes.Hum, mes.CO, mes.TubeHum, mes.TubeTemp))
	// writeApi.Flush()
	fmt.Println("getMessage")
}

func (app *application) mqqtHandler(c mqtt.Client) {
	c.Subscribe("mush/sensor-group/mesurments", 0, app.hadleIncomingMeasurements)
}
