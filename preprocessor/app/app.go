package app

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
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

type App struct {
	SQSClient     *SQSClient
	KinesisClient *KinesisClient
}

func (app *App) Init(creds Credentials) {
	var err error
	session := NewAWSSession(creds)

	app.SQSClient, err = NewSQSClient(session)
	if err != nil {
		log.Fatal("failed to create SQS client")
	}

	app.KinesisClient, err = NewKinesisClient(session)
	if err != nil {
		log.Fatal("failed to create Kinesis client")
	}
}

// TODO: testing purpose only
// Please use SQS lambda triggers
func (app *App) PollSQS(QueueUrl string, messageBatchLimit int64) error {
	ticker := time.NewTicker(time.Duration(10 * time.Second))
	for {
		select {
		case <-ticker.C:
			r := &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(QueueUrl),
				MaxNumberOfMessages: aws.Int64(messageBatchLimit),
				VisibilityTimeout:   aws.Int64(30),
				WaitTimeSeconds:     aws.Int64(20),
			}
			response, err := app.SQSClient.Read(r)
			if err != nil {
				log.Println(err)
				return err
			}

			for _, message := range response.Messages {
				messageBytes, err := base64.RawStdEncoding.DecodeString(*message.Body)
				if err != nil {
					log.Println(err)
					return err
				}
				var sqsMessage SQSMessage
				if err := json.Unmarshal(messageBytes, &sqsMessage); err != nil {
					log.Println(err)
					return err
				}

				id, err := uuid.NewUUID()
				if err != nil {
					log.Println(err)
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
					log.Println(err)
					return err
				}
				krWriteOut, err := app.KinesisClient.Write(&kinesis.PutRecordInput{
					Data:         krBytes,
					StreamName:   aws.String(os.Getenv("stream_name")),
					PartitionKey: aws.String("key1"),
				})
				if err != nil {
					log.Println(err)
					return err
				}
				log.Printf("Kinesis write output:%v", krWriteOut)

				sqsDeleteOut, err := app.SQSClient.Delete(&sqs.DeleteMessageInput{
					QueueUrl:      &QueueUrl,
					ReceiptHandle: message.ReceiptHandle,
				})
				if err != nil {
					log.Println(err)
					return err
				}
				log.Printf("SQS delete output:%v", sqsDeleteOut)
			}
		}
	}
}

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
