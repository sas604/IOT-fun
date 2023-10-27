package data

import (
	"errors"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Measurements MeasurementsModel
}

func NewModels(influxDB influxdb2.Client) Models {
	return Models{
		Measurements: MeasurementsModel{influxDb: influxDB},
	}
}
