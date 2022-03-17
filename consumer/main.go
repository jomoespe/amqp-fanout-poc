package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

const (
	exchangeName          = "poc.messages"
	alternateExchangeName = exchangeName + "-alternate"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Error to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// First create the alternate exchange, the alternate queue and bind it
	err = ch.ExchangeDeclare(
		alternateExchangeName, // name
		"fanout",              // kind
		true,                  // durable
		false,                 // auto-delete
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	failOnError(err, fmt.Sprintf("Failed to declare the alternate exchange %s", alternateExchangeName))

	alternateQueue, err := ch.QueueDeclare(
		fmt.Sprintf("%s@%s", "alternate", alternateExchangeName), // name. Queue name empty will generate a unique name which will be returned in the Name field
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare the alternate queue")

	err = ch.QueueBind(
		alternateQueue.Name,   // name
		"",                    // key
		alternateExchangeName, // exchange
		false,                 // no-wait
		nil,                   // arguments
	)
	failOnError(err, fmt.Sprintf("Failed to bind the alternate queue %s to alternate exchange %s", alternateQueue.Name, alternateExchangeName))

	// Then create the messaging/events exchange, the alternate queue and bind it
	// When declaring the exchange  we configure the alternate exchange
	err = ch.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // kind
		true,         // durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		map[string]interface{}{ // arguments
			"alternate-exchange": alternateExchangeName,
		},
	)
	failOnError(err, fmt.Sprintf("Failed to declare the exchange %s", exchangeName))

	q, err := ch.QueueDeclare(
		queueName(), // name. Queue name empty will generate a unique name which will be returned in the Name field
		false,       // durable
		true,        // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,       // name
		"",           // key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, fmt.Sprintf("Failed to bind the queue %s to exchange %s", q.Name, exchangeName))

	// Then we consume all messages from alternate queue
	pendingMessages, err := ch.Consume(
		alternateQueue.Name, // queue
		"",                  // consumer. Is unique and scoped for all consumers on this channel. An empty string will cause the library to generate a unique identity
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	failOnError(err, fmt.Sprintf("Failed to register a consumer to alternate queue %s", q.Name))

	log.Printf(" [*] Waiting for unprocessed messages.")
	go func() {
		for d := range pendingMessages {
			log.Printf("Processing pending message: %s", d.Body)
		}
		err := ch.QueueUnbind(
			alternateQueue.Name,   // name
			"",                    // key
			alternateExchangeName, // exchange
			nil,                   // args
		)
		if err != nil {
			log.Printf("Error unbinding from queue %s", alternateQueue.Name)
		}
	}()

	// At last, we consume all messages from queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer. Is unique and scoped for all consumers on this channel. An empty string will cause the library to generate a unique identity
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, fmt.Sprintf("Failed to register a consumer to queue %s", q.Name))

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("Queue name: %s bound to %s", q.Name, exchangeName)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func queueName() string {
	return fmt.Sprintf("%s@%s", uuid.NewString(), exchangeName)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
