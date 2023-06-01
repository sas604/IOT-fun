package db

import (
	"context"
	"errors"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// DB is a global variable to hold db connection
var DB influxdb2.Client

// ConnectDB opens a connection to the database
func ConnectToInfluxDb() error {
	dbToken := os.Getenv("INFLUXDB_TOKEN")
	if dbToken == "" {
		return errors.New("INFLUXDB_TOKEN must be set")
	}

	dbURL := os.Getenv("INFLUXDB_URL")
	if dbURL == "" {
		return errors.New("INFLUXDB_URL must be set")
	}
	client := influxdb2.NewClient(dbURL, dbToken)

	//validate client connection health
	_, err := client.Health(context.Background())

	DB = client

	return err
}
