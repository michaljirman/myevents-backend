package main

import (
	"fmt"

	"github.com/michaljirman/myevents-backend/src/bookingservice/listener"
	"github.com/michaljirman/myevents-backend/src/bookingservice/rest"
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
	var eventListener msgqueue.EventListener
	var eventEmitter msgqueue.EventEmitter

	//extract configuration
	config, _ := configuration.ExtractConfiguration()
	fmt.Println(config)
	switch config.MessageBrokerType {
	case "amqp":
		conn, err := amqp.Dial(config.AMQPMessageBroker)
		panicIfErr(err)

		eventListener, err = msgqueue_amqp.NewAMQPEventListener(conn, "events", "booking")
		panicIfErr(err)

		eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
		panicIfErr(err)
	case "kafka":
		conf := sarama.NewConfig()
		conf.Producer.Return.Successes = true
		conn, err := sarama.NewClient(config.KafkaMessageBrokers, conf)
		panicIfErr(err)

		eventListener, err = kafka.NewKafkaEventListener(conn, []int32{})
		panicIfErr(err)

		eventEmitter, err = kafka.NewKafkaEventEmitter(conn)
		panicIfErr(err)
	default:
		panic("Bad message broker type: " + config.MessageBrokerType)
	}

	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)

	processor := listener.EventProcessor{eventListener, dbhandler}
	go processor.ProcessEvents()

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
