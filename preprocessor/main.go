package main

import "os"

func main() {
	endpoint := os.Getenv("end_point")
	creds := Credentials{
		AccessKey:       os.Getenv("access_key_id"),
		SecretAccessKey: os.Getenv("secret_access_key"),
		Region:          os.Getenv("region"),
		EndPoint:        &endpoint,
	}
	app := new(App)
	app.Init(creds)
}
