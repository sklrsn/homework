package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	log.Println("************** Processing init starts**************")
	fmt.Println(os.Getenv("AWS_ACCESS_KEY_ID"))
	log.Println(os.Getenv("AWS_SECRET_ACCESS_KEY"))
	log.Println("************** Processing init starts**************")
}

func preprocess(ctx context.Context) error {
	log.Println("************** Processing starts**************")
	log.Println("************** Processing ends**************")
	return nil
}

func main() {
	lambda.Start(preprocess)
}
