package plug

import (
	"context"
	"encoding/json"

	"time"

	"github.com/looplab/fsm"
	"github.com/sas604/IOT-fun/server/db"
	mqttclient "github.com/sas604/IOT-fun/server/mqttClient"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Plug struct {
	BaseTopic string
	Switches  map[string]Switch
	PlugState string
}

type Switch struct {
	id  string
	FSM *fsm.FSM
}

func NewSwitch(id string, state string) Switch {
	s := Switch{
		id: id,
	}
	s.FSM = fsm.NewFSM(
		"initial",
		fsm.Events{
			{Name: "off", Src: []string{"on", "initial"}, Dst: "off"},
			{Name: "on", Src: []string{"off", "initial"}, Dst: "on"},
		},
		fsm.Callbacks{
			"enter_state": func(ctx context.Context, e *fsm.Event) { s.enterState(e) },
		},
	)
	s.FSM.Event(context.Background(), state)
	return s
}
func (s *Switch) enterState(e *fsm.Event) {
	writeApi := db.DB.WriteAPIBlocking("me", "iot-fun")

	p := influxdb2.NewPoint("sht-31", map[string]string{"sensor": "sht-31", "outlet": s.id}, map[string]interface{}{"state": e.Dst}, time.Now())

	writeApi.WritePoint(context.Background(), p)
	m, err := json.Marshal(map[string]string{"switch": s.id, "value": e.Dst})
	if err != nil {
		return
	}
	mqttclient.Client.Publish("mush/switch-group/set/"+s.id, 0, true, m)

}

func (p *Plug) SetSwitchStates(id string, state string) {
	if p.Switches[id].FSM.Cannot(state) {
		return
	}

	p.Switches[id].FSM.Event(context.Background(), state)

}

func NewPlug(s map[string]string) Plug {
	p := Plug{
		Switches:  make(map[string]Switch),
		BaseTopic: "mush/switch-group",
		PlugState: "offline",
	}

	for k, v := range s {
		p.Switches[k] = NewSwitch(k, v)
	}

	mqttclient.Client.Subscribe(p.BaseTopic+"/controllerStatus", 0, func(c MQTT.Client, m MQTT.Message) {
		p.PlugState = string(m.Payload())
	})
	return p
}
