package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"

	"github.com/hernanrocha/fin-chat/rabbit"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
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
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	rb := rabbit.NewRabbitChannel(ch)
	err = rb.QueueDeclare()
	failOnError(err, "Failed to declare a queue")

	msgs, err := rb.Consume()
	failOnError(err, "Failed to register a consumer")

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		str, _ := Process(string(d.Body))
		fmt.Println(str)
	}
}
