package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thiagonache/aws-marketplace-saas-integration/entitlement"
)

func main() {
	dynamoDBTableName := os.Getenv("AMSI_SUBSCRIBERS_TABLE_NAME")
	if dynamoDBTableName == "" {
		log.Fatal("Missing required environment variable AMSI_SUBSCRIBERS_TABLE_NAME")
	}
	e, err := entitlement.New(dynamoDBTableName)
	if err != nil {
		log.Fatal(err)
	}
	lambda.Start(e.HandleEntitlementMessage)
}
