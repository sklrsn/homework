package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/sklrsn/homework/preprocessor/app"
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

func main() {
	setEnv()
	endpoint := os.Getenv("end_point")
	creds := app.Credentials{
		AccessKey:       os.Getenv("access_key_id"),
		SecretAccessKey: os.Getenv("secret_access_key"),
		Region:          os.Getenv("region"),
		EndPoint:        &endpoint,
	}
	pp := new(app.App)
	pp.Init(creds)

	sqsbatchSize, err := strconv.Atoi(os.Getenv("sqs_messages_batch"))
	if err != nil {
		log.Fatalf("incorrect input for batch size :%v", err)
	}
	sqsPollInterval, err := strconv.Atoi(os.Getenv("sqs_poll_interval"))
	if err != nil {
		log.Fatalf("incorrect input for SQS poll interval :%v", err)
	}
	err = pp.PollSQS(os.Getenv("queue_url"), os.Getenv("stream_name"),
		int64(sqsbatchSize), sqsPollInterval)
	if err != nil {
		log.Printf("failed to process messages from SQS:%v", err)
	}

	done := make(chan bool)
	cleanUp(done)
	<-done
}

func setEnv() {
	os.Setenv("access_key_id", "AKIA2RVR24VPMP788U3S")
	os.Setenv("secret_access_key", "mDOWS+uN7dogVkHDaTuHaoyQ29Ju7pJmsvrrug8o")
	os.Setenv("region", "eu-west-1")
	os.Setenv("end_point", "http://localhost:4566")
	os.Setenv("queue_url", "http://localhost:4566/000000000000/submissions")
	os.Setenv("stream_name", "events")
	os.Setenv("sqs_messages_batch", "10")
	os.Setenv("sqs_poll_interval", "10") // in seconds
}
