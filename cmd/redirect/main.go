package main

import (
	"log"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thiagonache/aws-marketplace-saas-integration/redirect"
)

func main() {
	redirectLocation := os.Getenv("AMSI_REDIRECT_LOCATION")
	if redirectLocation == "" {
		log.Fatal("Missing required environment variable AMSI_REDIRECT_LOCATION")
	}
	u, err := url.Parse(redirectLocation)
	if err != nil {
		log.Fatal(err)
	}
	r := redirect.Redirect{Location: *u}
	lambda.Start(r.HandleRedirect)
}
