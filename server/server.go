package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sas604/IOT-fun/server/db"
	m "github.com/sas604/IOT-fun/server/mqttClient"
	"github.com/sas604/IOT-fun/server/plug"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func monitorMeasurement(d influxdb2.Client, c MQTT.Client) {

	tempTarget, err := strconv.ParseFloat(os.Getenv("TARGET_TEMP"), 64)
	if err != nil {
		fmt.Printf("Handle conversion error")
	}
	humTarget, err := strconv.ParseFloat(os.Getenv("TARGET_HUM"), 64)
	if err != nil {
		fmt.Printf("Handle conversion error")
	}
	//coTarget, err := strconv.ParseFloat(os.Getenv("TARGET_CO"), 64)

	if err != nil {
		fmt.Printf("Handle conversion error")
	}

	p := plug.NewPlug(map[string]string{"hum": "off", "heat": "off", "fan": "off", "light": "off"})
	for range time.Tick(time.Second * 60) {
		queryAPI := d.QueryAPI("me")
		fluxQuery := fmt.Sprintf(`from(bucket: "iot-fun")
		|> range(start: -1m)
		|> filter(fn: (r) => r["_measurement"] == "sht-31")
		|> filter(fn: (r) => r["_field"] == "co" or r["_field"] == "hum" or r["_field"] == "temp")
		|> median()`)
		result, err := queryAPI.Query(context.Background(), fluxQuery)
		if err != nil {
			// handle error
			fmt.Println("Error in DB query : ")
			fmt.Println(err)
		}
		for result.Next() {

			switch field := result.Record().Field(); field {
			case "hum":
				v := result.Record().Value().(float64)
				if v > humTarget && p.Switches["hum"].FSM.Current() == "on" {
					p.SetSwitchStates("hum", "off")
				}
				if v < humTarget && p.Switches["hum"].FSM.Current() == "off" {
					p.SetSwitchStates("hum", "on")
				}

			case "temp":

				v := result.Record().Value().(float64)
				if v > tempTarget && p.Switches["heat"].FSM.Current() == "on" {
					p.SetSwitchStates("heat", "off")
				}
				if v < tempTarget && p.Switches["heat"].FSM.Current() == "off" {
					p.SetSwitchStates("heat", "on")
				}
			// case "co":
			// 	v := result.Record().Value().(float64)
			// 	if v > coTarget {
			// 		fmt.Println(v, coTarget)
			// 	}

			default:
				//fmt.Printf("unrecognized field %s.\n", field)

			}

		}

	}
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal("Error loading .env file")
	}
	err = db.ConnectToInfluxDb()
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal("Error loading .env file")
	}
	m.ConectToMqttBroker()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go monitorMeasurement(db.DB, m.Client)
	<-c

}
