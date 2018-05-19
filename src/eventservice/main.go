package main

import (
	"fmt"

	"github.com/michaljirman/myevents-backend/src/eventservice/rest"
	"github.com/michaljirman/myevents-backend/src/lib/configuration"
	"github.com/michaljirman/myevents-backend/src/lib/msgqueue"
	msgqueue_amqp "github.com/michaljirman/myevents-backend/src/lib/msgqueue/amqp"
	"github.com/michaljirman/myevents-backend/src/lib/msgqueue/kafka"
	"github.com/michaljirman/myevents-backend/src/lib/persistence/dblayer"

	"github.com/Shopify/sarama"
	"github.com/streadway/amqp"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var eventEmitter msgqueue.EventEmitter

	// extract configuration
	config, _ := configuration.ExtractConfiguration()
	fmt.Println(config)

	switch config.MessageBrokerType {
	case "amqp":
		conn, err := amqp.Dial(config.AMQPMessageBroker)
		panicIfErr(err)

		eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
		panicIfErr(err)
	case "kafka":
		conf := sarama.NewConfig()
		conf.Producer.Return.Successes = true
		conn, err := sarama.NewClient(config.KafkaMessageBrokers, conf)
		panicIfErr(err)

		eventEmitter, err = kafka.NewKafkaEventEmitter(conn)
		panicIfErr(err)

	default:
		panic("Bad message broker type: " + config.MessageBrokerType)
	}

	fmt.Println("Connecting to database")
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)

	fmt.Println("Serving API")
	//RESTful API start
	httpErrChan, httptlsErrChan := rest.ServeAPI(config.RestfulEndpoint,
		config.RestfulTLSEndpoint, dbhandler, eventEmitter)
	select {
	case err := <-httpErrChan:
		panic(err)
	case err := <-httptlsErrChan:
		panic(err)
	}
}
