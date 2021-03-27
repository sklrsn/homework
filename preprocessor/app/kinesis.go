package app

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

type KinesisClient struct {
	kinesisClient *kinesis.Kinesis
}

func NewKinesisClient(session *session.Session) (*KinesisClient, error) {
	kinesisClient := new(KinesisClient)
	return kinesisClient, kinesisClient.initialize(session)
}

func (kc *KinesisClient) initialize(session *session.Session) error {
	kc.kinesisClient = kinesis.New(session)
	if kc.kinesisClient == nil {
		return errors.New("failed to initialize Kinesis handle")
	}
	return nil
}

func (kc *KinesisClient) Read() error {
	return errors.New("Not implemented")
}

func (kc *KinesisClient) Write(record *kinesis.PutRecordInput) (*kinesis.PutRecordOutput, error) {
	return kc.kinesisClient.PutRecord(record)
}

func (kc *KinesisClient) WriteMultipleRecords(records *kinesis.PutRecordsInput) (*kinesis.PutRecordsOutput, error) {
	return kc.kinesisClient.PutRecords(records)
}

func (kc *KinesisClient) Delete() error {
	return errors.New("Not implemented")
}
