package configuration

import (
	"fmt"
	"os"
	"strings"

	"github.com/michaljirman/myevents-backend/src/lib/persistence/dblayer"
)

var (
	DBTypeDefault              = dblayer.DBTYPE("mongodb")
	DBConnectionDefault        = "mongodb://127.0.0.1"
	RestfulEPDefault           = "localhost:8181"
	RestfulTLSEPDefault        = "localhost:9191"
	MessageBrokerTypeDefault   = "amqp"
	AMQPMessageBrokerDefault   = "amqp://guest:guest@localhost:5672"
	KafkaMessageBrokersDefault = []string{"localhost:9092"}
)

type ServiceConfig struct {
	Databasetype        dblayer.DBTYPE `json:"databasetype"`
	DBConnection        string         `json:"dbconnection"`
	RestfulEndpoint     string         `json:"restfulapi_endpoint"`
	RestfulTLSEndpoint  string         `json:"restfulapitls_endpoint"`
	MessageBrokerType   string         `json:"message_broker_type"`
	AMQPMessageBroker   string         `json:"amqp_message_broker"`
	KafkaMessageBrokers []string       `json:"kafka_message_brokers"`
}

func (config ServiceConfig) String() string {
	return fmt.Sprintf("ServiceConfig: \n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s",
		config.Databasetype,
		config.DBConnection,
		config.RestfulEndpoint,
		config.RestfulTLSEndpoint,
		config.MessageBrokerType,
		config.AMQPMessageBroker,
		config.KafkaMessageBrokers)
}

func ExtractConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		RestfulTLSEPDefault,
		MessageBrokerTypeDefault,
		AMQPMessageBrokerDefault,
		KafkaMessageBrokersDefault,
	}

	if v := os.Getenv("LISTEN_URL"); v != "" {
		conf.RestfulEndpoint = v
	}

	if v := os.Getenv("LISTEN_URL_TLS"); v != "" {
		conf.RestfulTLSEndpoint = v
	}

	if v := os.Getenv("MONGO_URL"); v != "" {
		conf.Databasetype = "mongodb"
		conf.DBConnection = v
	}

	if v := os.Getenv("AMQP_BROKER_URL"); v != "" {
		conf.MessageBrokerType = "amqp"
		conf.AMQPMessageBroker = v
	} else if v := os.Getenv("KAFKA_BROKER_URLS"); v != "" {
		conf.MessageBrokerType = "kafka"
		conf.KafkaMessageBrokers = strings.Split(v, ",")
	}
	return conf, nil
}
