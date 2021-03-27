package app

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type App struct {
	SQSClient     *SQSClient
	KinesisClient *KinesisClient
}

func (a *App) Init(creds Credentials) {
	session := NewAWSSession(creds)

	a.SQSClient = NewSQSClient(session)
	if a.SQSClient == nil {
		log.Fatal("failed to create SQS client")
	}

	a.KinesisClient = NewKinesisClient(session)
	if a.SQSClient == nil {
		log.Fatal("failed to create Kinesis client")
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
		},
	})
	if err != nil {
		log.Fatalf("incorrect credentials%v", err)
	}
	return session
}
