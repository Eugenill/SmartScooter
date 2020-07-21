package mqtt

import (
	"fmt"
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
	MQTTCLientID = "server"
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

func ListenToTopic(mqttConf MQTTConfig, topic string, client mqtt.Client) {
	client.Subscribe(mqttConf.Pretopic+topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
}
