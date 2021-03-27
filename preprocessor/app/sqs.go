package app

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSClient struct {
	sqsClient *sqs.SQS
}

func NewSQSClient(session *session.Session) *SQSClient {
	sqsClient := new(SQSClient)
	sqsClient.initialize(session)
	return sqsClient
}

func (s *SQSClient) initialize(session *session.Session) {
	s.sqsClient = sqs.New(session)
	if s.sqsClient == nil {
		log.Fatal("failed to initialize SQS")
	}
}

func (s *SQSClient) Read() {

}

func (s *SQSClient) Write() {

}

func (s *SQSClient) Delete() {

}
