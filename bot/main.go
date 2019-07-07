package main

import (
	"log"
	"os"

	"github.com/streadway/amqp"

	"github.com/hernanrocha/fin-chat/bot/stooq"
	"github.com/hernanrocha/fin-chat/messenger"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	rabbconn := getEnv("RABBIT_CONNECTION", "amqp://rabbitmq:rabbitmq@localhost:5672/")
	conn, err := amqp.Dial(rabbconn)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	log.Println("Starting to process messages")
	rb := messenger.NewRabbitMessenger(ch)
	err = rb.StartHandler(stooq.StooqHandler)
	failOnError(err, "Failed rabbit consumer")
}
