package mqtt_sub

import (
	"fmt"
	"github.com/Eugenill/SmartScooter/api_rest/handlers"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"net/url"
	"time"
)

const (
	MQTTUsername = "eugeni.llagostera@gmail.com"
	MQTTHost     = "maqiatto.com"
	MQTTPort     = ":1883"
	MQTTPassw    = "asdf1234"
	MQTTPreTopic = "eugeni.llagostera@gmail.com/"
	MQTTCLientID = "SmartScooter server"
)

type MQTTConfig struct {
	Host     string
	Port     string
	User     *url.Userinfo
	Pretopic string
	ClientID string
}

func ConnectToBroker(clientId string, mqttConf MQTTConfig) mqtt.Client {
	opts := CreateClientOptions(clientId, mqttConf)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func CreateClientOptions(clientId string, mqttConf MQTTConfig) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", mqttConf.Host+mqttConf.Port))
	opts.SetUsername(mqttConf.User.Username())
	password, _ := mqttConf.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func ListenToTopics(mqttConf MQTTConfig, topics []string, client mqtt.Client) {
	for _, topic := range topics {
		switch topic {
		case "timer":
			go client.Subscribe(mqttConf.Pretopic+topic, 0, handlers.LogDetection)
		case "RP1_detection":
			go client.Subscribe(mqttConf.Pretopic+topic, 0, handlers.SaveDetection)
		case "detection_example":
			go client.Subscribe(mqttConf.Pretopic+topic, 0, handlers.LogDetection)
		}
	}
	log.Printf("Topics added correctly")

}
