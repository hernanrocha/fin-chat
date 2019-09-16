package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	
	"github.com/hernanrocha/fin-chat/bot/stooq"
)

func main() {
	lambda.Start(stooq.StooqHandler)
}