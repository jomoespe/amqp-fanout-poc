package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

const (
	exchangeName = "direct"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: produce MSG")
		os.Exit(1)
	}
	msg := os.Args[1]

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Error to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.Publish(
		exchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
