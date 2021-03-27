package app

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSClient struct {
	sqs *sqs.SQS
}

func NewSQSClient(session *session.Session) *SQSClient {
	sqsClient := new(SQSClient)
	sqsClient.initialize(session)
	return sqsClient
}

func (s *SQSClient) initialize(session *session.Session) {
	s.sqs = sqs.New(session)
	if s.sqs == nil {
		log.Fatal("failed to initialize SQS Handle")
	}
}

func (s *SQSClient) Read(message *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return s.sqs.ReceiveMessage(message)
}

func (s *SQSClient) Write(message *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return s.sqs.SendMessage(message)
}

func (s *SQSClient) Delete() error {
	return errors.New("Not Implemented")
}
