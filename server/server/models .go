package server

import (
	"context"
	"fmt"

	"github.com/sas604/IOT-fun/server/db"
)

type SwitchWithState struct {
	State       string `json:"state"`
	Measurment  string `json:"measurment"`
	Value       int    `json:"value"`
	AutoControl bool   `json:"autoControl"`
	Target      int    `json:"target"`
	Unit        string `json:"unit"`
}

var measurmentToMesMap = map[string]string{
	"heat":  "Temperature",
	"temp":  "Temperature",
	"hum":   "Humidity",
	"co":    "CO2",
	"light": "Light",
}

func GetAllSwitchesWithState() ([]SwitchWithState, error) {
	queryAPI := db.DB.QueryAPI("me")
	fluxQuery := fmt.Sprintf(`from(bucket: "iot-fun")
  		|> range(start: 0)
  		|> filter(fn: (r) => r["_measurement"] == "plug")
  		|> filter(fn: (r) => r["_field"] == "co" or r["_field"] == "hum" or r["_field"] == "state" or r["_field"] == "temp")
		|> last()`)

	result, err := queryAPI.Query(context.Background(), fluxQuery)

	if err != nil {
		return nil, err
	}

	var res []SwitchWithState
	for result.Next() {
		mes, ok := measurmentToMesMap[result.Record().ValueByKey("outlet").(string)]
		if !ok {
			continue
		}
		res = append(res, SwitchWithState{Measurment: mes, State: result.Record().Value().(string)})

	}

	return res, nil

}
