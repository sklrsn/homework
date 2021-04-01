package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sklrsn/homework/preprocessor/app"
)

var (
	modeStandalone = "standalone"
)

func cleanUp(done chan bool) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		sig := <-ch
		log.Printf("Received signal %v", sig)

		log.Println("exiting ...")
		done <- true
	}()
}

var (
	Preprocessor = new(app.App)
)

func init() {
	endpoint := os.Getenv("end_point")
	creds := app.Credentials{
		AccessKey:       os.Getenv("access_key_id"),
		SecretAccessKey: os.Getenv("secret_access_key"),
		Region:          os.Getenv("region"),
		EndPoint:        &endpoint,
	}
	Preprocessor.Init(creds)

	batch, err := strconv.Atoi(os.Getenv("sqs_messages_batch"))
	if err != nil {
		log.Fatalf("incorrect batch size:%v", err)
	}
	interval, err := strconv.Atoi(os.Getenv("sqs_poll_interval"))
	if err != nil {
		log.Fatalf("incorrect poll interval:%v", err)
	}
	timeout, err := strconv.Atoi(os.Getenv("visibility_timeout"))
	if err != nil {
		log.Fatalf("incorrect visibility timeout:%v", err)
	}
	Preprocessor.Params = &app.Params{
		BatchSize:    batch,
		PollInterval: interval,
		Queue:        os.Getenv("queue_url"),
		Stream:       os.Getenv("stream_name"),
		Timeout:      timeout,
	}
}

func handler(ctx context.Context, event events.SQSEvent) error {
	messages := make([]*sqs.Message, 0)
	for _, record := range event.Records {
		message := sqs.Message{
			MessageId:     &record.MessageId,
			Body:          &record.Body,
			ReceiptHandle: &record.ReceiptHandle,
		}
		messages = append(messages, &message)
	}

	return Preprocessor.ProcessSQSEvent(Preprocessor.Params, messages)
}

func main() {
	if os.Getenv("mode") == modeStandalone {
		err := Preprocessor.PollSQS(Preprocessor.Params)
		if err != nil {
			log.Printf("failed to process messages from SQS:%v", err)
		}

		done := make(chan bool)
		cleanUp(done)
		<-done
	}

	lambda.Start(handler)
}
