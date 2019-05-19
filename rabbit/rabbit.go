package rabbit

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitChannel interface {
	QueueDeclare() error
	Publish(key, message string) error
	Consume() (<-chan amqp.Delivery, error)
	Handle(fn func(string) (string, error)) error
}

func NewRabbitChannel(ch *amqp.Channel) RabbitChannel {
	return &rabbitChannel{
		Ch: ch,
	}
}

type rabbitChannel struct {
	Ch   *amqp.Channel
	req  amqp.Queue
	resp amqp.Queue
}

func (r *rabbitChannel) QueueDeclare() error {
	// Declare request queue
	req, err := r.Ch.QueueDeclare(
		"req_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	r.req = req

	if err != nil {
		return err
	}

	resp, err := r.Ch.QueueDeclare(
		"resp_queue", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	r.resp = resp

	return err
}

func (r *rabbitChannel) Publish(key, message string) error {
	return r.Ch.Publish(
		"",         // exchange
		r.req.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: key,
			Body:          []byte(message),
			ReplyTo:       r.resp.Name,
		})
}

func (r *rabbitChannel) Handle(fn func(string) (string, error)) error {
	msgs, err := r.Ch.Consume(
		r.req.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return err
	}
	for msg := range msgs {
		// Process request and generate response
		log.Printf("Received a message: %s", msg.Body)
		resp, _ := fn(string(msg.Body))
		fmt.Println(resp)

		r.Ch.Publish(
			"",          // exchange
			msg.ReplyTo, // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				Body:          []byte(resp),
				CorrelationId: msg.CorrelationId,
			})
	}
	return nil
}

func (r *rabbitChannel) Consume() (<-chan amqp.Delivery, error) {
	return r.Ch.Consume(
		r.resp.Name, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
}
