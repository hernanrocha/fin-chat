package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hernanrocha/fin-chat/rabbit"
	"github.com/streadway/amqp"
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

func Process(s string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://stooq.com/q/l/?s=%s.us&f=sd2t2ohlcv&h&e=csv", s))
	reader := csv.NewReader(bufio.NewReader(resp.Body))
	_, err = reader.Read()
	if err != nil {
		log.Printf("Error reading header: %s \n", err)
		return fmt.Sprintf("Error obtaining info for %s", s), err
	}
	row, err := reader.Read()
	if err != nil || len(row) <= 4 {
		log.Printf("Error reading row: %s \n", err)
		return fmt.Sprintf("Error obtaining info for %s", s), err
	}
	return fmt.Sprintf("%s quote is $%s per share", s, row[3]), nil
}

func main() {
	rabbconn := getEnv("RABBIT_CONNECTION", "amqp://rabbitmq:rabbitmq@localhost:5672/")
	conn, err := amqp.Dial(rabbconn)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	rb := rabbit.NewRabbitChannel(ch)
	err = rb.QueueDeclare()
	failOnError(err, "Failed to declare a queue")

	fmt.Println("Starting to process messages")
	rb.Handle(Process)
}
