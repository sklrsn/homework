package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sklrsn/homework/preprocessor/app"
)

const (
	QueueUrl = "http://localhost:4566/000000000000/submissions"
)

func main() {
	setEnv()
	fmt.Println("**************")
	fmt.Println(os.Getenv("end_point"))
	fmt.Println(os.Getenv("access_key_id"))
	fmt.Println(os.Getenv("secret_access_key"))
	fmt.Println(os.Getenv("region"))
	fmt.Println("**************")
	endpoint := os.Getenv("end_point")
	creds := app.Credentials{
		AccessKey:       os.Getenv("access_key_id"),
		SecretAccessKey: os.Getenv("secret_access_key"),
		Region:          os.Getenv("region"),
		EndPoint:        &endpoint,
	}
	app := new(app.App)
	app.Init(creds)

	s := &sqs.SendMessageInput{
		MessageBody:  aws.String("message body"),
		QueueUrl:     aws.String(QueueUrl),
		DelaySeconds: aws.Int64(3),
	}
	response, err := app.SQSClient.Write(s)
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
	receive_resp, err := app.SQSClient.Read(r)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("[Receive message] \n%v \n\n", receive_resp)
}

func setEnv() {
	os.Setenv("access_key_id", "AKIA2RVR24VPMP788U3S")
	os.Setenv("secret_access_key", "mDOWS+uN7dogVkHDaTuHaoyQ29Ju7pJmsvrrug8o")
	os.Setenv("region", "eu-west-1")
	os.Setenv("end_point", "http://localhost:4566")
	os.Setenv("queue_url", "http://localhost:4566/000000000000/submissions")
}
