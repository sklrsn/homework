package app

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

const (
	UTC = "2006-01-02T15:04:05.999999Z"
)

type Credentials struct {
	AccessKey       string  `json:"access_key"`
	SecretAccessKey string  `json:"secret_access_key"`
	SessionToken    string  `json:"session_token"`
	Region          string  `json:"region"`
	EndPoint        *string `json:"end_point"`
}

func (c Credentials) String() string {
	return ""
}

func NewAWSSession(creds Credentials) *session.Session {
	cred := credentials.NewStaticCredentials(creds.AccessKey, creds.SecretAccessKey,
		creds.SessionToken)
	session, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigDisable,
		Config: aws.Config{
			Region:      aws.String(creds.Region),
			Credentials: cred,
			Endpoint:    creds.EndPoint,
			MaxRetries:  aws.Int(5),
		},
	})
	if err != nil {
		log.Fatalf("incorrect credentials%v", err)
	}
	return session
}

type Params struct {
	BatchSize    int
	PollInterval int
	Queue        string
	Stream       string
}
type App struct {
	SQSClient     *SQSClient
	KinesisClient *KinesisClient
	Params        *Params
}

func (app *App) Init(creds Credentials) {
	var err error
	session := NewAWSSession(creds)

	app.SQSClient, err = NewSQSClient(session)
	if err != nil {
		log.Fatalf("failed to create SQS client:%v", err)
	}

	app.KinesisClient, err = NewKinesisClient(session)
	if err != nil {
		log.Fatalf("failed to create Kinesis client:%v", err)
	}
}

// TODO: testing purpose only
// Please use SQS lambda triggers
func (app *App) PollSQS(QueueUrl, stream string, batchSize int64, interval int) error {
	if interval <= 0 {
		interval = 10
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer func() {
					wg.Done()
				}()

				response, err := app.readSQSMessage(QueueUrl, batchSize)
				if err != nil {
					log.Println(err)
					return
				}

				for _, message := range response.Messages {
					if sqsMessage, err := app.DecodeSQSMessage(*message.Body); err == nil {
						if err := app.WriteToKinesisStream(stream, *sqsMessage); err != nil {
							log.Println(err)
							return
						}
					}
					if err := app.DeleteSQSMessage(QueueUrl, message.ReceiptHandle); err != nil {
						log.Println(err)
						return
					}
				}
			}(&wg)

			wg.Wait()
		}
	}
}

func (app *App) readSQSMessage(QueueUrl string, batchSize int64) (*sqs.ReceiveMessageOutput, error) {
	r := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(QueueUrl),
		MaxNumberOfMessages: aws.Int64(batchSize),
		VisibilityTimeout:   aws.Int64(30),
		WaitTimeSeconds:     aws.Int64(20),
	}
	response, err := app.SQSClient.Read(r)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (app *App) DecodeSQSMessage(body string) (*SQSMessage, error) {
	messageBytes, err := base64.RawStdEncoding.DecodeString(body)
	if err != nil {
		return nil, err
	}
	var sqsMessage SQSMessage
	if err := json.Unmarshal(messageBytes, &sqsMessage); err != nil {
		return nil, err
	}
	return &sqsMessage, err
}

func (app *App) DeleteSQSMessage(QueueUrl string, receiptHandle *string) error {
	sqsDeleteOut, err := app.SQSClient.Delete(&sqs.DeleteMessageInput{
		QueueUrl:      &QueueUrl,
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		return err
	}
	log.Printf("SQS delete output:%v", sqsDeleteOut)

	return nil
}

func (app *App) WriteToKinesisStream(stream string, sqsMessage SQSMessage) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	kinesisRecord := KinesisRecord{
		RecordID:           id.String(),
		DeviceID:           sqsMessage.DeviceID,
		Processes:          sqsMessage.Events.Processes,
		NetworkConnections: sqsMessage.Events.NetworkConnections,
		Created:            time.Now().UTC().Format(UTC),
	}
	krBytes, err := json.Marshal(kinesisRecord)
	if err != nil {
		return err
	}
	krWriteOut, err := app.KinesisClient.Write(&kinesis.PutRecordInput{
		Data:         krBytes,
		StreamName:   aws.String(stream),
		PartitionKey: aws.String("key1"),
	})
	if err != nil {
		return err
	}
	log.Printf("Kinesis write output:%v", krWriteOut)

	return nil
}
