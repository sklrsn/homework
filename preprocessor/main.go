package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func preprocess(ctx context.Context, event events.SQSEvent) error {
	log.Println("************** SQS event Processing starts**************")
	for _, message := range event.Records {
		log.Println(message.MessageId)
		log.Println(message.Body)
	}
	log.Println("************** SQS event Processing ends**************")
	return nil
}

func main() {
	lambda.Start(preprocess)
}
