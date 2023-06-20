package entitlement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/marketplaceentitlementservice"
)

type Entitlement struct {
	GetEntitlements       func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error)
	SubscribersTableName  string
	UpdateItemWithContext func(context.Context, *dynamodb.UpdateItemInput, ...request.Option) (*dynamodb.UpdateItemOutput, error)
}

func New(subscribersTableName string) (Entitlement, error) {
	if subscribersTableName == "" {
		return Entitlement{}, errors.New("cannot be empty")
	}
	sess := session.Must(session.NewSession())
	awsConfig := aws.NewConfig()
	m := marketplaceentitlementservice.New(sess, awsConfig)
	d := dynamodb.New(sess, awsConfig)
	return Entitlement{
		GetEntitlements:       m.GetEntitlements,
		SubscribersTableName:  subscribersTableName,
		UpdateItemWithContext: d.UpdateItemWithContext,
	}, nil
}

func (e Entitlement) HandleEntitlementMessage(ctx context.Context, input *events.SQSEvent) error {
	type entitlementMessage struct {
		Action             string `json:"action"`
		CustomerIdentifier string `json:"customer-identifier"`
		ProductCode        string `json:"product-code"`
	}
	type sqsEntitlementMessage struct {
		Type    string
		Message entitlementMessage
	}
	if len(input.Records) != 1 {
		return fmt.Errorf("wrong number of records in the event (%d). Please, configure lambda trigger batch size to one", len(input.Records))
	}
	event := input.Records[0]
	msg := &sqsEntitlementMessage{}
	err := json.Unmarshal([]byte(event.Body), msg)
	if err != nil {
		return err
	}
	if msg.Message.Action != "entitlement-updated" {
		return fmt.Errorf("invalid action in message %q", event.Body)
	}
	entitlementOutput, err := e.GetEntitlements(&marketplaceentitlementservice.GetEntitlementsInput{
		Filter: map[string][]*string{
			"CUSTOMER_IDENTIFIER": {&msg.Message.CustomerIdentifier},
		},
		ProductCode: &msg.Message.ProductCode,
	})
	if err != nil {
		return err
	}
	expired := false
	if entitlementOutput.Entitlements[0].ExpirationDate.Before(time.Now()) {
		expired = true
	}
	entitlementData, err := json.Marshal(entitlementOutput)
	if err != nil {
		return err
	}
	_, err = e.UpdateItemWithContext(ctx, &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":e":  {S: aws.String(string(entitlementData))},
			":se": {BOOL: aws.Bool(expired)},
			":ss": {BOOL: aws.Bool(true)},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"customerIdentifier": {S: &msg.Message.CustomerIdentifier},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		TableName:        &e.SubscribersTableName,
		UpdateExpression: aws.String("set entitlement = :e, successfully_subscribed = :ss, subscription_expired = :se"),
	})
	if err != nil {
		return err
	}
	return nil
}
