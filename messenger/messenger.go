package messenger

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/streadway/amqp"
)

type BotMessage struct {
	RoomID  uint
	Message string
}

type BotCommandMessenger interface {
	// Publish a command request message
	Publish(roomID uint, message string) error
	// Start command request message handler
	// StartHandler(func(string) (string, error)) error
	// Start command response message consumer
	StartConsumer(func(BotMessage) error) error
}

func NewRabbitMessenger(ch *amqp.Channel) *rabbitCommandMessenger {
	return &rabbitCommandMessenger{
		Ch: ch,
	}
}

func NewSQSMessenger(snsSvc *sns.SNS, sqsSvc *sqs.SQS) *sqsCommandMessenger {
	return &sqsCommandMessenger{
		sqsSvc: sqsSvc,
		snsSvc: snsSvc,
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

type sqsCommandMessenger struct {
	snsSvc *sns.SNS
	sqsSvc *sqs.SQS
}

// Publish a command request message
func (s *sqsCommandMessenger) Publish(roomID uint, message string) error {
	req := &BotMessage{
		Message: message,
		RoomID:  roomID,
	}
	resStr, _ := json.Marshal(req)
	_, err := s.snsSvc.Publish(&sns.PublishInput{
		Message:  aws.String(string(resStr)),
		TopicArn: aws.String(os.Getenv("SQS_COMMANDS_REQUEST_URL")),
	})
	return err
}

// Start command response message consumer
func (s *sqsCommandMessenger) StartConsumer(fn func(BotMessage) error) error {
	log.Println("Starting command message response consumer")

	for {
		output, err := s.sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(os.Getenv("SQS_COMMANDS_RESPONSE_URL")),
			MaxNumberOfMessages: aws.Int64(1),
			WaitTimeSeconds:     aws.Int64(20),
		})

		if err != nil {
			log.Printf("Failed to fetch sqs message %v", err)
			return err
		}

		for _, msg := range output.Messages {
			var res BotMessage

			// Parse message
			if err := json.Unmarshal([]byte(*msg.Body), &res); err != nil {
				return err
			}

			// Process message
			if err := fn(res); err != nil {
				log.Printf("Error processing message: %s\n", err.Error())
			}

			// Delete message
			dmr := &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(os.Getenv("SQS_COMMANDS_RESPONSE_URL")),
				ReceiptHandle: msg.ReceiptHandle,
			}
			if _, err = s.sqsSvc.DeleteMessage(dmr); err != nil {
				log.Printf("Error deleting message from queue: %s\n", err.Error())
			}
		}
	}
}
