package rabbit

import (
	"github.com/streadway/amqp"
)

type RabbitChannel interface {
	QueueDeclare() error
	Publish(message string) error
	Consume() (<-chan amqp.Delivery, error)
}

func NewRabbitChannel(ch *amqp.Channel) RabbitChannel {
	return &rabbitChannel{
		Ch: ch,
	}
}

type rabbitChannel struct {
	Ch *amqp.Channel
	q  amqp.Queue
}

func (r *rabbitChannel) QueueDeclare() error {
	q, err := r.Ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	r.q = q
	return err
}

func (r *rabbitChannel) Publish(message string) error {
	return r.Ch.Publish(
		"",       // exchange
		r.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (r *rabbitChannel) Consume() (<-chan amqp.Delivery, error) {
	return r.Ch.Consume(
		r.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
}
