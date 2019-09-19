package stooq

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type BotMessage struct {
	RoomID  uint
	Message string
}

func StooqHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	if len(sqsEvent.Records) == 0 {
		return errors.New("No SQS message passed to function")
	}

	mySession := session.New()
	svc := sqs.New(mySession)

	for _, msg := range sqsEvent.Records {
		fmt.Printf("Got SQS message %q with body %q\n", msg.MessageId, msg.Body)
		var req BotMessage
		if err := json.Unmarshal([]byte(msg.Body), &req); err != nil {
			return err
		}

		res, err := Handle(req.Message)
		if err != nil {
			return err
		}

		resStr, _ := json.Marshal(&BotMessage{RoomID: req.RoomID, Message: res})
		_, err = svc.SendMessage(&sqs.SendMessageInput{
			MessageBody: aws.String(string(resStr)),
			QueueUrl:    aws.String(os.Getenv("SQS_COMMANDS_RESPONSE_URL")),
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func Handle(s string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://stooq.com/q/l/?s=%s.us&f=sd2t2ohlcv&h&e=csv", s))
	if err != nil {
		return fmt.Sprintf("Error getting HTTP response for %s", s), err
	}

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
	log.Printf("%s quote is $%s per share \n", s, row[3])
	return fmt.Sprintf("%s quote is $%s per share", s, row[3]), nil
}
