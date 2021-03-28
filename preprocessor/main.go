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
	pp = new(app.App)
)

func init() {
	endpoint := os.Getenv("end_point")
	creds := app.Credentials{
		AccessKey:       os.Getenv("access_key_id"),
		SecretAccessKey: os.Getenv("secret_access_key"),
		Region:          os.Getenv("region"),
		EndPoint:        &endpoint,
	}
	pp.Init(creds)

	batch, err := strconv.Atoi(os.Getenv("sqs_messages_batch"))
	if err != nil {
		log.Fatalf("incorrect batch size:%v", err)
	}
	interval, err := strconv.Atoi(os.Getenv("sqs_poll_interval"))
	if err != nil {
		log.Fatalf("incorrect poll interval:%v", err)
	}
	pp.Params = &app.Params{
		BatchSize:    batch,
		PollInterval: interval,
		Queue:        os.Getenv("queue_url"),
		Stream:       os.Getenv("stream_name"),
	}
}

func handler(ctx context.Context, events events.SQSEvent) error {
	for _, message := range events.Records {
		if sqsMessage, err := pp.DecodeSQSMessage(message.Body); err == nil {
			if err := pp.WriteToKinesisStream(pp.Params.Stream, *sqsMessage); err != nil {
				log.Println(err)
				continue
			}
		}
		if err := pp.DeleteSQSMessage(pp.Params.Queue, &message.ReceiptHandle); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func main() {
	if os.Getenv("mode") == modeStandalone {
		err := pp.PollSQS(pp.Params.Queue, pp.Params.Stream,
			int64(pp.Params.BatchSize), pp.Params.PollInterval)
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

	log.Fatalf("incorrect mode")
}

func setEnvironment() {
	os.Setenv("access_key_id", "AKIA2RVR24VPMP788U3S")
	os.Setenv("secret_access_key", "mDOWS+uN7dogVkHDaTuHaoyQ29Ju7pJmsvrrug8o")
	os.Setenv("region", "eu-west-1")
	os.Setenv("end_point", "http://localhost:4566")
	os.Setenv("queue_url", "http://localhost:4566/000000000000/submissions")
	os.Setenv("stream_name", "events")
	os.Setenv("sqs_messages_batch", "10")
	os.Setenv("sqs_poll_interval", "10") // in seconds
	os.Setenv("mode", "standalone")
}
