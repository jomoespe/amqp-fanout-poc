package main

import (
	"log"

	"github.com/streadway/amqp"
)

const exchange = "poc.messages"

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Error to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name. Queue name empty will generate a unique name which will be returned in the Name field
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,   // name
		"",       // key
		exchange, // exchange
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to bind to a queue "+q.Name)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer. Is unique and scoped for all consumers on this channel. An empty string will cause the library to generate a unique identity
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer to queue "+q.Name)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("Queue name: %s", q.Name)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
