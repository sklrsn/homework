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
	"github.com/sklrsn/homework/preprocessor/app"
)

var (
	modeStandalone = "standalone"
	modeServerless = "serverless"
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
	Preprocessor.Params = &app.Params{
		BatchSize:    batch,
		PollInterval: interval,
		Queue:        os.Getenv("queue_url"),
		Stream:       os.Getenv("stream_name"),
	}
}

func handler(ctx context.Context, events events.SQSEvent) error {
	for _, message := range events.Records {
		if sqsMessage, err := Preprocessor.DecodeSQSMessage(message.Body); err == nil {
			if err := Preprocessor.WriteToKinesisStream(Preprocessor.Params.Stream,
				*sqsMessage); err != nil {
				log.Println(err)
				continue
			}
		}
		if err := Preprocessor.DeleteSQSMessage(Preprocessor.Params.Queue,
			&message.ReceiptHandle); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func main() {
	if os.Getenv("mode") == modeStandalone {
		err := Preprocessor.PollSQS(Preprocessor.Params.Queue, Preprocessor.Params.Stream,
			int64(Preprocessor.Params.BatchSize), Preprocessor.Params.PollInterval)
		if err != nil {
			log.Printf("failed to process messages from SQS:%v", err)
		}

		done := make(chan bool)
		cleanUp(done)
		<-done
	}

	if os.Getenv("mode") == modeServerless {
		lambda.Start(handler)
	}

	log.Fatalf("incorrect mode:%v . supported modes -{'standalone', 'serverless'} ",
		os.Getenv("mode"))
}
