package data

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type MeasurementsModel struct {
	influxDb influxdb2.Client
}

type Measurements struct {
	Measurement string  `json:"measurment,omitempty"`
	Temp        float64 `json:"temp,omitempty"`
	Hum         float64 `json:"hum,omitempty"`
	CO          float64 `json:"co,omitempty"`
	TubeHum     float64 `json:"tube_hum,omitempty"`
	TubeTemp    float64 `json:"tube_temp,omitempty"`
}

func (m MeasurementsModel) Insert(mes *Measurements, org string, bucket string) error {
	writeAPI := m.influxDb.WriteAPIBlocking(org, bucket)
	p := fmt.Sprintf("%s,sensor=sht-31 temp=%f,hum=%f,co=%f,tube_hum=%f,tube_temp=%f", mes.Measurement, mes.Temp, mes.Hum, mes.CO, mes.TubeHum, mes.TubeTemp)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return writeAPI.WriteRecord(ctx, p)

}
