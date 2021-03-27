package main

import (
	"fmt"
	"os"

	"github.com/sklrsn/homework/preprocessor/app"
)

func main() {
	fmt.Println("**************")
	fmt.Println(os.Getenv("end_point"))
	fmt.Println(os.Getenv("access_key_id"))
	fmt.Println(os.Getenv("secret_access_key"))
	fmt.Println(os.Getenv("region"))
	fmt.Println("**************")
	endpoint := os.Getenv("end_point")
	creds := app.Credentials{
		AccessKey:       os.Getenv("access_key_id"),
		SecretAccessKey: os.Getenv("secret_access_key"),
		Region:          os.Getenv("region"),
		EndPoint:        &endpoint,
	}
	app := new(app.App)
	app.Init(creds)
	fmt.Println(app.SQSClient)
	fmt.Println(app.KinesisClient)
}
