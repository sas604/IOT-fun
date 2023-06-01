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

	MQTT "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type SensorData struct {
	Temp float64 `json:"temp"`
	Hum  float64 `json:"hum"`
	CO   float64 `json:"co"`
}

func monitorMeasurement(client influxdb2.Client, ch chan SensorData) {
	for range time.Tick(time.Second * 5) {
		queryAPI := client.QueryAPI("me")
		fluxQuery := fmt.Sprintf(`from(bucket: "iot-fun")
		|> range(start: -1h)
		|> filter(fn: (r) => r["_measurement"] == "sht-31")
		|> filter(fn: (r) => r["_field"] == "co" or r["_field"] == "hum" or r["_field"] == "temp")
		|> median()`)
		result, err := queryAPI.Query(context.Background(), fluxQuery)
		if err != nil {
			// handle error
			fmt.Println(err)
		}

		var resultPoints SensorData
		for result.Next() {
			switch field := result.Record().Field(); field {
			case "hum":
				resultPoints.Hum = result.Record().Value().(float64)
			case "temp":
				resultPoints.Temp = result.Record().Value().(float64)
			case "co":
				resultPoints.CO = result.Record().Value().(float64)
			default:
				fmt.Printf("unrecognized field %s.\n", field)

			}

		}
		fmt.Println("sending", resultPoints)
		ch <- resultPoints
	}
}

func handleMeasurementReadings(ch chan SensorData) {
	var err error
	tempTarget, err := strconv.ParseFloat(os.Getenv("TARGET_TEMP"), 64)
	// humTarget, err := strconv.ParseFloat(os.Getenv("TARGET_HUM"), 64)
	// coTarget, err := strconv.ParseFloat(os.Getenv("TARGET_CO"), 64)
	if err != nil {
		fmt.Printf("Handle cversion error")
	}

	for {
		v := <-ch
		if v.Temp > tempTarget {
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

// func connectToInfluxDb() (influxdb2.Client, error) {
// 	dbToken := os.Getenv("INFLUXDB_TOKEN")
// 	if dbToken == "" {
// 		return nil, errors.New("INFLUXDB_TOKEN must be set")
// 	}

// 	dbURL := os.Getenv("INFLUXDB_URL")
// 	if dbURL == "" {
// 		return nil, errors.New("INFLUXDB_URL must be set")
// 	}
// 	client := influxdb2.NewClient(dbURL, dbToken)

// 	//validate client connection health
// 	_, err := client.Health(context.Background())

// 	return client, err
// }

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
	mes := make(chan SensorData)
	go monitorMeasurement(db.DB, mes)
	go handleMeasurementReadings(mes)

	<-c

}
