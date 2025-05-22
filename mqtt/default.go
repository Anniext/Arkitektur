package mqtt

import (
	"time"
)

var (
	defaultMqtt *MQTTClient
)

func InitDefaultMqtt() error {
	cnf := GetDefaultMqttConfig()

	defaultMqtt = NewClient(ClientConfig{
		BrokerURL:            cnf.BrokerURL,
		ClientID:             cnf.ServerID,
		KeepAlive:            time.Duration(cnf.KeepAlive) * time.Second,
		QoS:                  cnf.QoS,
		Retain:               cnf.Retain,
		AutoReconnect:        cnf.AutoReconnect,
		ConnectTimeout:       time.Duration(cnf.ConnectTimeout) * time.Second,
		MaxReconnectInterval: time.Duration(cnf.MaxReconnectInterval) * time.Second,
	})

	if err := defaultMqtt.Connect(); err != nil {
		return err
	}

	return nil
}

func GetDefaultMqtt() *MQTTClient {
	return defaultMqtt
}
