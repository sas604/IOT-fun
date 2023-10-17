package data

import influxdb2 "github.com/influxdata/influxdb-client-go/v2"

type MeasurementsModel struct {
	influxDb *influxdb2.Client
}
