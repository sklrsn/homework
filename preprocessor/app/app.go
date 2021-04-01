package app

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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

var (
	ErrMissingEvents                      = errors.New("SQS event is missing events")
	ErrMissingProcessAndNetworkConnection = errors.New("SQS event is missing new_process and network_connection")
	ErrMissingDeviceID                    = errors.New("SQS event is missing device_id")
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
	Timeout      int
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

func (app *App) PollSQS(params *Params) error {
	interval := params.PollInterval
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

				response, err := app.ReadSQSMessage(params.Queue, int64(params.BatchSize), int64(params.Timeout))
				if err != nil {
					log.Println(err)
					return
				}

				if err := app.ProcessSQSEvent(params, response.Messages); err != nil {
					log.Println(err)
					return
				}

			}(&wg)

			wg.Wait()
		}
	}
}

func (app *App) ProcessSQSEvent(params *Params, messages []*sqs.Message) error {
	if sqsMessages, err := app.DecodeSQSMessages(messages); err == nil {
		if err := app.ValidateSQSMessages(sqsMessages); err == nil {
			err := app.WriteMessagesToKinesisStream(params.Stream, sqsMessages)
			if err != nil {
				return err
			}
		}
	} else {
		log.Println(err)
	}

	for _, message := range messages {
		err := app.DeleteSQSMessage(*message.MessageId, params.Queue, message.ReceiptHandle)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) ValidateSQSMessages(messages []SQSMessage) error {
	if len(messages) == 0 {
		return ErrMissingEvents
	}

	for _, message := range messages {
		if len(message.Events.Processes) == 0 && len(message.Events.NetworkConnections) == 0 {
			return ErrMissingProcessAndNetworkConnection
		}

		if len(message.DeviceID) == 0 {
			return ErrMissingDeviceID
		}
	}

	return nil
}

func (app *App) ReadSQSMessage(QueueUrl string, batchSize, timeout int64) (*sqs.ReceiveMessageOutput, error) {
	r := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(QueueUrl),
		MaxNumberOfMessages: aws.Int64(batchSize),
		VisibilityTimeout:   aws.Int64(timeout),
		WaitTimeSeconds:     aws.Int64(20),
	}
	response, err := app.SQSClient.Read(r)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (app *App) DecodeSQSMessages(messages []*sqs.Message) ([]SQSMessage, error) {
	response := make([]SQSMessage, 0)

	for _, message := range messages {
		messageBytes, err := base64.RawStdEncoding.DecodeString(*message.Body)
		if err != nil {
			return nil, err
		}

		var sqsMessage SQSMessage
		if err := json.Unmarshal(messageBytes, &sqsMessage); err != nil {
			return nil, err
		}
		response = append(response, sqsMessage)
	}

	return response, nil
}

func (app *App) DeleteSQSMessage(messageID, QueueUrl string, receiptHandle *string) error {
	_, err := app.SQSClient.Delete(&sqs.DeleteMessageInput{
		QueueUrl:      &QueueUrl,
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		return err
	}
	log.Printf("SQS removed messsageID:%v", messageID)

	return nil
}

func (app *App) WriteMessagesToKinesisStream(stream string, sqsMessage []SQSMessage) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	kinesisRecord := KinesisRecord{
		RecordID: id.String(),
		Data:     sqsMessage,
		Created:  time.Now().UTC().Format(UTC),
	}
	krBytes, err := json.Marshal(kinesisRecord)
	if err != nil {
		return err
	}
	log.Printf("record written to Kinesis:%v", string(krBytes))

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
