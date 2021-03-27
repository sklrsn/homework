package app

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSClient struct {
	sqs *sqs.SQS
}

func NewSQSClient(session *session.Session) (*SQSClient, error) {
	sqsClient := new(SQSClient)
	return sqsClient, sqsClient.initialize(session)
}

func (sc *SQSClient) initialize(session *session.Session) error {
	sc.sqs = sqs.New(session)
	if sc.sqs == nil {
		return errors.New("failed to initialize SQS")
	}
	return nil
}

func (sc *SQSClient) Read(message *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return sc.sqs.ReceiveMessage(message)
}

func (sc *SQSClient) Write(message *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return sc.sqs.SendMessage(message)
}

func (sc *SQSClient) Delete(deleteMessage *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return sc.sqs.DeleteMessage(deleteMessage)
}
