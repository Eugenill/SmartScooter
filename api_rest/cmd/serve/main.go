package main

import (
	"context"
	"database/sql"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/Eugenill/SmartScooter/api_rest/router"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
)

var mqttConfig = mqtt_sub.MQTTConfig{
	Host:     mqtt_sub.MQTTHost,
	Port:     mqtt_sub.MQTTPort,
	User:     url.UserPassword(mqtt_sub.MQTTUsername, mqtt_sub.MQTTPassw),
	Pretopic: mqtt_sub.MQTTPreTopic,
	ClientID: mqtt_sub.MQTTCLientID,
}

func main() {
	//1. Create Context
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()
	//2. Open DB?
	db, err := sql.Open("postgres", "host=localhost port=5432 dbname=smartscooter user=postgres password=postgres sslmode=disable")
	errors.Catch(err)
	log.Printf("Connected to the PostgreSQL, Host: localhost, Port: 5432, dbname: smartscooter, user: postgres, password: postgres")
	ctx = bunny.ContextWithDB(ctx, db)

	//3. Initialize MQTT
	client := mqtt_sub.ConnectToBroker(mqttConfig.ClientID, mqttConfig)
	log.Printf("Connected to the MQTTbroker, CLient:%s, Host:%s, Port%s", mqttConfig.ClientID, mqttConfig.Host, mqttConfig.Port)
	topics := []string{"timer", "detection_example", "RP1_detection"}
	mqtt_sub.ListenToTopics(mqttConfig, topics, client)
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
