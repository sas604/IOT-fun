package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"sync"

	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sas604/IOT-fun/server/internal/data"
)

type config struct {
	port     int
	env      string
	influxDB struct {
		influxToken string
		influxURL   string
		bucket      string
		org         string
	}

	mqtt struct {
		brokerAddr string
		password   string
		userName   string
	}

	// limiter struct {
	// 	rps     float64
	// 	burst   int
	// 	enabled bool
	// }
	// smtp struct {
	// 	host     string
	// 	port     int
	// 	username string
	// 	password string
	// 	sender   string
	// }
	// cors struct {
	// 	trustedOrigins []string
	// }
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
	// mailer mailer.Mailer
	wg sync.WaitGroup
}

func main() {

	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.influxDB.influxToken, "influx-token", os.Getenv("INFLUX_TOKEN"), "Influx DB Token")
	flag.StringVar(&cfg.influxDB.influxURL, "influx-url", "http://192.168.1.106:8086", "Influx url")
	flag.StringVar(&cfg.influxDB.bucket, "bucket", "monitoring", "Influx bucket name")
	flag.StringVar(&cfg.influxDB.org, "org", "me", "Influx org name")
	flag.StringVar(&cfg.mqtt.brokerAddr, "mqqt-url", "tcp://192.168.1.106:1883", "Mqtt broker url")
	flag.StringVar(&cfg.mqtt.password, "mqtt-password", os.Getenv("MQTT_PASS"), "Mqtt broker password")
	flag.StringVar(&cfg.mqtt.userName, "mqtt-user-name", os.Getenv("MQTT_USER_NAME"), "Mqtt broker user Name")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	influxClient, err := newInfluxClient(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer influxClient.Close()
	logger.Info("got influx client")

	db, err := sql.Open("sqlite3", "./sqlite.db")

	if err != nil {
		logger.Error("error conecting to db", "error", err.Error())
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		logger.Error("error conecting to db", "error", err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("connected to the db")

	var mqttClient mqtt.Client

	app := &application{
		config: cfg,
		logger: logger,
	}

	mqttClient, err = app.newMQTTClient(cfg)
	if err != nil {
		logger.Error("error setting up mqtt controler " + err.Error())
		os.Exit(1)
	}
	defer mqttClient.Disconnect(100)

	logger.Info("Conected to mqtt")

	app.models = data.NewModels(influxClient, db, mqttClient, app.config.influxDB.org, app.config.influxDB.bucket)

	go app.handlePeriodicTasks()
	err = app.listnAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func newInfluxClient(cfg config) (influxdb2.Client, error) {

	client := influxdb2.NewClient(cfg.influxDB.influxURL, cfg.influxDB.influxToken)

	//validate client connection health
	_, err := client.Health(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil

}

func (app *application) newMQTTClient(cfg config) (mqtt.Client, error) {
	options := mqtt.NewClientOptions().AddBroker(cfg.mqtt.brokerAddr)
	options.SetUsername(cfg.mqtt.userName)
	options.SetPassword(cfg.mqtt.password)
	options.SetOrderMatters(false)
	options.OnConnect = app.mqqtHandler
	client := mqtt.NewClient(options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}
