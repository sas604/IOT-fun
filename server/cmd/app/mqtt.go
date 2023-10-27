package main

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sas604/IOT-fun/server/internal/data"
)

func (app *application) hadleIncomingMeasurements(c mqtt.Client, m mqtt.Message) {
	var mes data.Measurements
	err := json.Unmarshal([]byte(m.Payload()), &mes)

	if err != nil {
		app.logger.Error(err.Error())
		return
	}
	err = app.models.Measurements.Insert(&mes, app.config.influxDB.org, app.config.influxDB.bucket)
	if err != nil {
		app.logger.Error("failed insert", err.Error())
		return
	}

}

func (app *application) mqqtHandler(c mqtt.Client) {
	c.Subscribe("mush/sensor-group/mesurments", 0, app.hadleIncomingMeasurements)
}
