package plug

import (
	"encoding/json"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Plug struct {
	BaseTopic string
	Switches  map[string]string
	db        influxdb2.Client
	client    MQTT.Client
	PlugState string
}

func (p *Plug) GetSwitchState(id string) {
	fmt.Println(p.Switches[id])
}

func (p *Plug) SetSwitchStates(id string, state string) {
	fmt.Println("id", id, state)
	if state != "off" && state != "on" {
		return
	}
	if p.Switches[id] == state {
		return
	}
	m, err := json.Marshal(map[string]string{"switch": id, "value": state})
	if err != nil {
		return
	}
	p.client.Publish("mush/switch-group/set/"+id, 0, true, m)
	p.Switches[id] = state
}

var onTransitionResult MQTT.MessageHandler = func(c MQTT.Client, m MQTT.Message) {
	fmt.Println(string(m.Payload()))
}

func NewPlug(db influxdb2.Client, c MQTT.Client, s map[string]string) Plug {
	p := Plug{
		Switches:  s,
		BaseTopic: "mush/switch-group",
		db:        db,
		client:    c,
		PlugState: "offline",
	}

	c.Subscribe(p.BaseTopic+"/controllerStatus", 0, func(c MQTT.Client, m MQTT.Message) {
		fmt.Println(string(m.Payload()))
		p.PlugState = string(m.Payload())
	})
	for id := range p.Switches {
		jm, _ := json.Marshal(map[string]string{"switch": id, "value": p.Switches[id]})
		c.Publish(p.BaseTopic+"/set/"+id, 0, true, jm)

	}

	c.Subscribe("mush/switch-group/transition", 0, onTransitionResult)
	return p
}
