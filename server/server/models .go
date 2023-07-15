package server

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/sas604/IOT-fun/server/db"
)

type SwitchWithState struct {
	State       string  `json:"state"`
	Measurment  string  `json:"measurment"`
	Value       float64 `json:"value"`
	AutoControl bool    `json:"autoControl"`
	Target      int64   `json:"target"`
	Unit        string  `json:"unit"`
}

var measurmentToMesMap = map[string]string{
	"temp":  "Temperature",
	"hum":   "Humidity",
	"co":    "CO2",
	"light": "Light",
}

var getAllStateQ = `from(bucket: "iot-fun")
|> range(start: 0)
|> filter(fn: (r) => r["_measurement"] == "sht-31")
|> filter(fn: (r) => r["_field"] == "state")
|> last()`

func GetAllSwitchesWithState() ([]SwitchWithState, error) {
	queryAPI := db.DB.QueryAPI("me")
	fluxQuery := getAllStateQ
	result, err := queryAPI.Query(context.Background(), fluxQuery)

	if err != nil {
		return nil, err
	}
	var res []SwitchWithState
	for result.Next() {
		key := result.Record().ValueByKey("outlet").(string)
		m, ok := measurmentToMesMap[key]
		fmt.Println(m, key)
		if !ok {
			continue
		}

		q := fmt.Sprintf(`from(bucket: "iot-fun")
		|> range(start: 0)
		|> filter(fn: (r) => r["_measurement"] == "sht-31")
		|> filter(fn: (r) => r["_field"] == "%s")
		|> last()`, key)

		r, err := queryAPI.Query(context.Background(), q)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		s := SwitchWithState{
			State:       result.Record().Value().(string),
			Measurment:  m,
			AutoControl: true,
		}

		switch key {
		case "fan":
			s.Unit = "ppm"
		case "hum":
			s.Unit = "%"
			h, err := strconv.ParseInt(os.Getenv("TARGET_HUM"), 0, 0)
			if err != nil {
				s.Target = 0
			}
			s.Target = h
		case "temp":
			s.Unit = "C"
			t, err := strconv.ParseInt(os.Getenv("TARGET_TEMP"), 0, 0)
			if err != nil {
				s.Target = 0
			}
			s.Target = t
		}

		if r.Next() {
			s.Value = r.Record().Value().(float64)
		}

		res = append(res, s)

	}

	return res, nil

}
