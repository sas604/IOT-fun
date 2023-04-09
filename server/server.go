package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type SensorData struct {
	name     string
	location string
	hum      string
	temp     string
}

func writeToDb(client influxdb2.Client, s SensorData) {
	// get non-blocking write client
	writeApi := client.WriteAPI("me", "iot-fun")
	writeApi.WriteRecord(fmt.Sprintf("sht-31,controller=main,sensor=%s,location=%s temp=%s,hum=%s", s.name, s.location, s.temp, s.hum))
	writeApi.Flush()
}

var msgH MQTT.MessageHandler = func(c MQTT.Client, m MQTT.Message) {

	values := strings.Split(string(m.Payload()), ",")
	data := SensorData{
		name:     "sht-31",
		location: "any",
		temp:     values[0],
		hum:      values[1],
	}

	writeToDb(db, data)

}
var db influxdb2.Client

func connectToInfluxDb() (influxdb2.Client, error) {
	dbToken := os.Getenv("INFLUXDB_TOKEN")
	if dbToken == "" {
		return nil, errors.New("INFLUXDB_TOKEN must be set")
	}

	dbURL := os.Getenv("INFLUXDB_URL")
	if dbURL == "" {
		return nil, errors.New("INFLUXDB_URL must be set")
	}
	client := influxdb2.NewClient(dbURL, dbToken)

	//validate client connection health
	_, err := client.Health(context.Background())

	return client, err
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal("Error loading .env file")
	}
	db, err = connectToInfluxDb()
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal("Error loading .env file")
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	options := MQTT.NewClientOptions().AddBroker("tcp://192.168.1.106:1883")
	options.SetUsername("pi")
	options.SetPassword("boopyou")
	options.SetOrderMatters(false)
	options.OnConnect = func(c MQTT.Client) {
		c.Subscribe("test/temperature", 0, msgH)
		c.Subscribe("test/humidity", 0, msgH)
	}
	client := MQTT.NewClient(options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server \n")
	}
	<-c

}
