package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type config struct {
	port int
	env  string
	// db   struct {
	// 	dsn          string
	// 	maxOpenConns int
	// 	maxIdleConns int
	// 	maxIdleTime  time.Duration
	// }
	influxDB struct {
		influxToken string
		influxURL   string
		bucket      string
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
	// models data.Models
	// mailer mailer.Mailer
	// wg     sync.WaitGroup
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.influxDB.influxToken, "influx-token", "", "Influx DB Token")
	flag.StringVar(&cfg.influxDB.influxURL, "influx-url", "http://192.168.1.106:8086", "Influx url")
	flag.StringVar(&cfg.influxDB.bucket, "bucket", "monitoring", "Influx bucket name")
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
	logger.Info("get influx client")
	app := &application{
		config: cfg,
		logger: logger,
	}

	mqttClient, err := app.newMQTTClient(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer mqttClient.Disconnect(100)

	fmt.Println(mqttClient.IsConnected())
	// app := &application{
	// 	config: cfg,
	// 	logger: loger,
	// }
	for {
		time.Sleep(10 * time.Second)
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