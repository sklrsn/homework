package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sklrsn/homework/preprocessor/app"
)

const (
	QueueUrl   = "http://localhost:4566/000000000000/submissions"
	StreamName = "events"
)

func cleanUp() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		sig := <-ch
		log.Printf("Received signal %v", sig)

		log.Println("exiting ...")
		os.Exit(0)
	}()
}

func main() {
	setEnv()
	fmt.Println("**************")
	fmt.Println(os.Getenv("end_point"))
	fmt.Println(os.Getenv("access_key_id"))
	fmt.Println(os.Getenv("secret_access_key"))
	fmt.Println(os.Getenv("region"))
	fmt.Println(os.Getenv("stream_name"))
	fmt.Println("**************")

	endpoint := os.Getenv("end_point")
	creds := app.Credentials{
		AccessKey:       os.Getenv("access_key_id"),
		SecretAccessKey: os.Getenv("secret_access_key"),
		Region:          os.Getenv("region"),
		EndPoint:        &endpoint,
	}
	pp := new(app.App)
	pp.Init(creds)

	s := &sqs.SendMessageInput{
		MessageBody:  aws.String("message body"),
		QueueUrl:     aws.String(QueueUrl),
		DelaySeconds: aws.Int64(3),
	}
	response, err := pp.SQSClient.Write(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(response)

	r := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(QueueUrl),
		MaxNumberOfMessages: aws.Int64(3),
		VisibilityTimeout:   aws.Int64(30),
		WaitTimeSeconds:     aws.Int64(20),
	}
	receive_resp, err := pp.SQSClient.Read(r)
	if err != nil {
		log.Println(err)
	}

	for _, message := range receive_resp.Messages {
		messageBytes, _ := base64.RawStdEncoding.DecodeString(*message.Body)
		fmt.Println(string(messageBytes))
	}

	putOutput, err := pp.KinesisClient.Write(&kinesis.PutRecordInput{
		Data:         []byte("hoge"),
		StreamName:   aws.String(os.Getenv("stream_name")),
		PartitionKey: aws.String("key1"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(putOutput)
}

func setEnv() {
	os.Setenv("access_key_id", "AKIA2RVR24VPMP788U3S")
	os.Setenv("secret_access_key", "mDOWS+uN7dogVkHDaTuHaoyQ29Ju7pJmsvrrug8o")
	os.Setenv("region", "eu-west-1")
	os.Setenv("end_point", "http://localhost:4566")
	os.Setenv("queue_url", "http://localhost:4566/000000000000/submissions")
	os.Setenv("stream_name", "events")
}
