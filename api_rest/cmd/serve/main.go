package main

import (
	"context"
	"database/sql"
	"github.com/Eugenill/SmartScooter/api_rest/mqtt"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/Eugenill/SmartScooter/api_rest/router"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
)

var mqttConfig = mqtt.MQTTConfig{
	Host:     mqtt.MQTTHost,
	Port:     mqtt.MQTTPort,
	User:     url.UserPassword(mqtt.MQTTUsername, mqtt.MQTTPassw),
	Pretopic: mqtt.MQTTPreTopic,
	ClientID: mqtt.MQTTCLientID,
}

func main() {
	//1. Create Context
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()
	//2. Open DB?
	db, err := sql.Open("postgres", "host=localhost port=5432 dbname=smartscooter user=postgres password=postgres sslmode=disable")
	errors.Catch(err)
	ctx = bunny.ContextWithDB(ctx, db)

	//3. Initialize MQTT
	client := mqtt.ConnectToBroker(mqttConfig.ClientID, mqttConfig)
	log.Printf("Connected to the MQTTbroker, Host:%s, Port%s", mqttConfig.Host, mqttConfig.Port)
	topics := []string{"timer", "detection"}
	for _, topic := range topics {
		go mqtt.ListenToTopic(mqttConfig, topic, client)
	}
	//mqtt_client.PublishTimer("timer", mqttConfig)

	//4. Init router
	errL := http.ListenAndServe("localhost:1234", router.SetRouter(mqttConfig))
	errors.Catch(errL)

	// Waiting for an OS signal cancellation
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	//log

	// Shutdown the servers

	//log
}
