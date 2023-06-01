package plug

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Plug struct {
	State     string
	BaseTopic string
	Switches  map[string]string
	db        influxdb2.Client
	client    MQTT.Client
}

func (p *Plug) Off(id string) {
	fmt.Println(p.State, id)
}

func NewPlug(db influxdb2.Client, c MQTT.Client) Plug {
	p := Plug{
		State:     "off",
		BaseTopic: "mush/switch-group",
		db:        db,
		client:    c,
	}
	return p
}
