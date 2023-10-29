package data

import (
	"database/sql"
	"errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Measurements MeasurementsModel
	Switches     SwitchModel
	Automations  AutomationModel
}

func NewModels(influxDB influxdb2.Client, db *sql.DB, MQTT mqtt.Client) Models {
	return Models{
		Measurements: MeasurementsModel{influxDb: influxDB},
		Switches:     SwitchModel{DB: db, MQTT: MQTT},
		Automations:  AutomationModel{DB: db},
	}
}
