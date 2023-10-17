package data

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Models struct {
	Measurements MeasurementsModel
}

func NewModels(influxDB *influxdb2.Client) Models {
	return Models{
		Measurements: MeasurementsModel{influxDb: *influxDB},
	}
}
