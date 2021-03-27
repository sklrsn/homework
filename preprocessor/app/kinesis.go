package app

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

type KinesisClient struct {
	kinesisClient *kinesis.Kinesis
}

func NewKinesisClient(session *session.Session) *KinesisClient {
	kinesisClient := new(KinesisClient)
	kinesisClient.initialize(session)
	return kinesisClient
}

func (s *KinesisClient) initialize(session *session.Session) {
	s.kinesisClient = kinesis.New(session)
	if s.kinesisClient == nil {
		log.Fatal("failed to initialize SQS")
	}
}

func (s *KinesisClient) Read() error {
	return errors.New("Not implemented")
}

func (s *KinesisClient) Write() {

}

func (s *KinesisClient) Delete() error {
	return errors.New("Not implemented")
}
