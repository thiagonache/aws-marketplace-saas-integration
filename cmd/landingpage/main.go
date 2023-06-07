package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/marketplacemetering"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/thiagonache/aws-marketplace-saas-integration/landingpage"
)

func main() {
	entitlementQueueURL := os.Getenv("AMSI_ENTITLEMENT_QUEUE_URL")
	if entitlementQueueURL == "" {
		log.Fatal("Missing required environment variable AMSI_ENTITLEMENT_QUEUE_URL")
	}
	dynamoDBTableName := os.Getenv("AMSI_SUBSCRIBERS_TABLE_NAME")
	if dynamoDBTableName == "" {
		log.Fatal("Missing required environment variable AMSI_SUBSCRIBERS_TABLE_NAME")
	}
	sess := session.Must(session.NewSession())
	awsConfig := aws.NewConfig()
	l, err := landingpage.New(entitlementQueueURL, dynamoDBTableName)
	if err != nil {
		log.Fatal(err)
	}
	mktplaceMeeting := marketplacemetering.New(sess, awsConfig)
	l.ResolveCustomerWithContext = mktplaceMeeting.ResolveCustomerWithContext
	sqs := sqs.New(sess, awsConfig)
	l.SendMessageWithContext = sqs.SendMessageWithContext
	dynamo := dynamodb.New(sess, awsConfig)
	l.PutItemWithContext = dynamo.PutItemWithContext
	lambda.Start(l.HandleLandingPage)
}
