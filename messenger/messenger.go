package messenger

import (
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

type BotMessage struct {
	RoomID  uint
	Message string
}

type BotCommandMessenger interface {
	// Publish a command request message
	Publish(key, message string) error
	// Start command request message handler
	StartHandler(func(string) (string, error)) error
	// Start command response message consumer
	StartConsumer(func(BotMessage) error) error
}

func NewRabbitMessenger(ch *amqp.Channel) *rabbitCommandMessenger {
	return &rabbitCommandMessenger{
		Ch: ch,
	}
}

type rabbitCommandMessenger struct {
	Ch   *amqp.Channel
	req  amqp.Queue
	resp amqp.Queue
}

func (r *rabbitCommandMessenger) queueDeclare() error {
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

func (r *rabbitCommandMessenger) Publish(key, message string) error {
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

func (r *rabbitCommandMessenger) StartHandler(fn func(string) (string, error)) error {
	if err := r.queueDeclare(); err != nil {
		return err
	}

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
		log.Printf("Received command response: %s", msg.Body)
		resp, _ := fn(string(msg.Body))
		log.Println(resp)

		err := r.Ch.Publish(
			"",          // exchange
			msg.ReplyTo, // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				Body:          []byte(resp),
				CorrelationId: msg.CorrelationId,
			})

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *rabbitCommandMessenger) StartConsumer(fn func(BotMessage) error) error {
	if err := r.queueDeclare(); err != nil {
		return err
	}

	msgs, err := r.Ch.Consume(
		r.resp.Name, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}

	log.Printf("Waiting for messages from rabbit...\n")

	for d := range msgs {
		log.Printf("Received a message: %s\n", d.Body)
		roomID, _ := strconv.Atoi(d.CorrelationId)
		msg := BotMessage{
			RoomID:  uint(roomID),
			Message: string(d.Body),
		}

		if err := fn(msg); err != nil {
			log.Printf("Error processing message: %s\n", err.Error())
		}
	}

	return nil
}
