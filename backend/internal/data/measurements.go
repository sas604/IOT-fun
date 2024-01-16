package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type MeasurementsModel struct {
	influxDb influxdb2.Client
	DB       *sql.DB
	org      string
	bucket   string
}

type Measurements struct {
	Measurement string  `json:"measurment,omitempty"` //TODO fix typo in sensor
	Temp        float64 `json:"temp,omitempty"`
	Hum         float64 `json:"hum,omitempty"`
	CO          float64 `json:"co,omitempty"`
	TubeHum     float64 `json:"tube_hum,omitempty"`
	TubeTemp    float64 `json:"tube_temp,omitempty"`
}

func (m MeasurementsModel) Insert(mes *Measurements, org string, bucket string) error {
	writeAPI := m.influxDb.WriteAPIBlocking(org, bucket)
	p := fmt.Sprintf("%s,sensor=sht-31 temp=%f,hum=%f,co=%f,tube_hum=%f,tube_temp=%f", "farm", mes.Temp, mes.Hum, mes.CO, mes.TubeHum, mes.TubeTemp)

	writeAPI.WriteRecord(context.Background(), p)

	return nil
}

func (m MeasurementsModel) GetMeasurementMap() (map[string]float64, error) {
	queryAPI := m.influxDb.QueryAPI(m.org)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		from(bucket: "monitoring")
  		|> range(start: 0)
  		|> filter(fn: (r) => r["_measurement"] == "farm")
  		|> last()`

	rows, err := queryAPI.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	ms := map[string]float64{}
	for rows.Next() {
		ms[rows.Record().Field()] = (rows.Record().Value()).(float64)
	}

	return ms, nil
}
