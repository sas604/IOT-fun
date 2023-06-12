package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sas604/IOT-fun/server/db"
	"github.com/sas604/IOT-fun/server/plug"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type SensorData struct {
	Temp float64 `json:"temp"`
	Hum  float64 `json:"hum"`
	CO   float64 `json:"co"`
}

func monitorMeasurement(d influxdb2.Client, c MQTT.Client) {

	tempTarget, err := strconv.ParseFloat(os.Getenv("TARGET_TEMP"), 64)
	if err != nil {
		fmt.Printf("Handle conversion error")
	}
	humTarget, err := strconv.ParseFloat(os.Getenv("TARGET_HUM"), 64)
	if err != nil {
		fmt.Printf("Handle conversion error")
	}
	coTarget, err := strconv.ParseFloat(os.Getenv("TARGET_CO"), 64)

	if err != nil {
		fmt.Printf("Handle conversion error")
	}

	p := plug.NewPlug(d, c, map[string]string{"hum": "off", "heat": "off", "fan": "off", "light": "off"})
	for range time.Tick(time.Second * 5) {
		queryAPI := d.QueryAPI("me")
		fluxQuery := fmt.Sprintf(`from(bucket: "iot-fun")
		|> range(start: -1m)
		|> filter(fn: (r) => r["_measurement"] == "sht-31")
		|> filter(fn: (r) => r["_field"] == "co" or r["_field"] == "hum" or r["_field"] == "temp")
		|> median()`)
		result, err := queryAPI.Query(context.Background(), fluxQuery)
		if err != nil {
			// handle error
			fmt.Println(err)
		}
		for result.Next() {
			switch field := result.Record().Field(); field {
			case "hum":
				v := result.Record().Value().(float64)
				if v > humTarget && p.Switches["hum"] == "on" {
					fmt.Println("set hum to off")
					p.SetSwitchStates("hum", "off")
				}
				if v < humTarget && p.Switches["hum"] == "off" {
					fmt.Println("set hum to on")
					p.SetSwitchStates("hum", "on")
				}

			case "temp":

				v := result.Record().Value().(float64)
				fmt.Println("temp is ", v, "State is", p.Switches["heat"])
				if v > tempTarget && p.Switches["heat"] == "on" {
					fmt.Println("set heat to off ")
					p.SetSwitchStates("heat", "off")
				}
				if v < tempTarget && p.Switches["heat"] == "off" {
					fmt.Println("set heat to on ")
					p.SetSwitchStates("heat", "on")
				}
			case "co":
				v := result.Record().Value().(float64)
				if v > coTarget {
					fmt.Println(v, coTarget)
				}

			default:
				fmt.Printf("unrecognized field %s.\n", field)

			}

		}

	}
}

func writeToDb(client influxdb2.Client, s SensorData) {
	// get non-blocking write client
	writeApi := client.WriteAPI("me", "iot-fun")
	writeApi.WriteRecord(fmt.Sprintf("sht-31,controller=main,sensor=sht-31,location=%s temp=%f,hum=%f,co=%f", "any", s.Temp, s.Hum, s.CO))
	writeApi.Flush()
}

var msgH MQTT.MessageHandler = func(c MQTT.Client, m MQTT.Message) {
	var mes SensorData
	err := json.Unmarshal([]byte(m.Payload()), &mes)
	if err != nil {
		fmt.Print(err.Error())

	}

	writeToDb(db.DB, mes)

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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	options := MQTT.NewClientOptions().AddBroker("tcp://192.168.1.106:1883")
	options.SetUsername("pi")
	options.SetPassword("boopyou")
	options.SetOrderMatters(false)
	options.OnConnect = func(c MQTT.Client) {
		c.Subscribe("mush/sensor-group/mesurments", 0, msgH)
	}
	client := MQTT.NewClient(options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server \n")
	}
	go monitorMeasurement(db.DB, client)
	<-c

}
