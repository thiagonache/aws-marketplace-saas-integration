package entitlement_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/marketplaceentitlementservice"
	"github.com/google/go-cmp/cmp"
	"github.com/thiagonache/aws-marketplace-saas-integration/entitlement"
)

var (
	defaultSQSEvent = &events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: `{
						"type": "Notification",
						"message" : {
							"action" : "entitlement-updated",
							"customer-identifier": "customerIdentifier",
							"product-code" : "productCode"
						}
					}`,
			},
		},
	}
)

func TestHandleEntitlementMessage_ErrorsGivenNoRecord(t *testing.T) {
	t.Parallel()
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		&events.SQSEvent{
			Records: []events.SQSMessage{},
		},
	)
	if err == nil {
		t.Fatal("want error but got nil")
	}
}

func TestHandleEntitlementMessage_ErrorsGivenMoreThanOneRecord(t *testing.T) {
	t.Parallel()
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		&events.SQSEvent{
			Records: []events.SQSMessage{
				{
					Body: `{
						"type": "Notification",
						"message" : {
							"action" : "entitlement-updated",
							"customer-identifier": "customerIdentifier",
							"product-code" : "productCode"
						}
					}`,
				},
				{
					Body: `{
						"type": "Notification",
						"message" : {
							"action" : "entitlement-updated",
							"customer-identifier": "customerIdentifier",
							"product-code" : "productCode2"
						}
					}`,
				},
			},
		},
	)
	if err == nil {
		t.Fatal("want error but got nil")
	}
}

func TestHandleEntitlementMessage_ErrorsIfActionIsUnexpected(t *testing.T) {
	t.Parallel()
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		&events.SQSEvent{
			Records: []events.SQSMessage{
				{
					Body: `{
						"type": "Notification",
						"message" : {
							"action" : "bogus",
							"customer-identifier": "customerIdentifier",
							"product-code" : "productCode"
						}
					}`,
				},
			},
		},
	)
	if err == nil {
		t.Fatal("want error but got nil")
	}
}

func TestNew_SetsGetEntitlementsByDefault(t *testing.T) {
	t.Parallel()
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	if e.GetEntitlements == nil {
		t.Fatal("want GetEntitlements to be set but got nil")
	}
}

func TestHandleEntitlementMessage_CallsGetEntitlements(t *testing.T) {
	t.Parallel()
	called := false
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		called = true
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		defaultSQSEvent,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("GetEntitlements not called")
	}
}

func TestHandleEntitlementMessage_SetsProperProductCodeInGetEntitlementsAPICall(t *testing.T) {
	t.Parallel()
	want := "productCode"
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		got := *input.ProductCode
		if want != got {
			t.Fatalf("want product code %q, got %q", want, got)
		}
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		defaultSQSEvent,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleEntitlementMessage_SetsProperFilterCustomerIdentifierInGetEntitlementsAPICall(t *testing.T) {
	t.Parallel()
	want := "customerIdentifier"
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		if len(input.Filter["CUSTOMER_IDENTIFIER"]) < 1 {
			t.Fatal("Filter CUSTOMER_IDENTIFIER not passed")
		}
		got := *input.Filter["CUSTOMER_IDENTIFIER"][0]
		if want != got {
			t.Fatalf("want filter CUSTOMER_IDENTIFIER %q, got %q", want, got)
		}
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		defaultSQSEvent,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNew_SetsUpdateItemWithContextByDefault(t *testing.T) {
	t.Parallel()
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	if e.UpdateItemWithContext == nil {
		t.Fatal("want UpdateItemWithContext to be set but got nil")
	}
}

func TestHandleEntitlementMessage_CallsUpdateItemWithContext(t *testing.T) {
	t.Parallel()
	called := false
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		called = true
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		defaultSQSEvent,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("UpdateItemWithContext not called")
	}
}

func TestNew_SetsSubscribersTableNameByDefault(t *testing.T) {
	t.Parallel()
	want := "myDynamoDBTableName"
	e, err := entitlement.New("myDynamoDBTableName")
	if err != nil {
		t.Fatal(err)
	}
	got := e.SubscribersTableName
	if want != got {
		t.Fatalf("want subscribers table name %q, got %q", want, got)
	}
}

func TestHandleEntitlementMessage_SetsProperTableNameInUpdateItemWithContextAPICall(t *testing.T) {
	t.Parallel()
	want := "MyDynamoDBTableName"
	e, err := entitlement.New("MyDynamoDBTableName")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		got := *input.TableName
		if want != got {
			t.Fatalf("want DynamoDB table name %q, got %q", want, got)
		}
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		defaultSQSEvent,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleEntitlementMessage_SetsProperKeyInUpdateItemWithContextAPICall(t *testing.T) {
	t.Parallel()
	want := &dynamodb.AttributeValue{
		S: aws.String("customerIdentifier"),
	}
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		got := input.Key["customerIdentifier"]
		if !cmp.Equal(want, got) {
			t.Fatal(cmp.Diff(want, got))
		}
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		defaultSQSEvent,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleEntitlementMessage_SetsProperUpdateExpressionInUpdateItemWithContextAPICall(t *testing.T) {
	t.Parallel()
	want := "set entitlement = :e, successfully_subscribed = :ss, subscription_expired = :se"
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		got := *input.UpdateExpression
		if !cmp.Equal(want, got) {
			t.Fatal(cmp.Diff(want, got))
		}
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		defaultSQSEvent,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleEntitlementMessage_SetsProperReturnValuesInUpdateItemWithContextAPICall(t *testing.T) {
	t.Parallel()
	want := "UPDATED_NEW"
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		got := *input.ReturnValues
		if !cmp.Equal(want, got) {
			t.Fatal(cmp.Diff(want, got))
		}
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		defaultSQSEvent,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleEntitlementMessage_SetsProperSSExpressionAttributeValueInUpdateItemWithContextAPICall(t *testing.T) {
	t.Parallel()
	want := &dynamodb.AttributeValue{
		BOOL: aws.Bool(true),
	}
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		got := input.ExpressionAttributeValues[":ss"]
		if !cmp.Equal(want, got) {
			t.Fatal(cmp.Diff(want, got))
		}
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		defaultSQSEvent,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleEntitlementMessage_SetsProperSEExpressionAttributeValueInUpdateItemWithContextAPICall(t *testing.T) {
	t.Parallel()
	want := &dynamodb.AttributeValue{
		BOOL: aws.Bool(false),
	}
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{ExpirationDate: aws.Time(time.Now().Add(time.Hour))},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		got := input.ExpressionAttributeValues[":se"]
		if !cmp.Equal(want, got) {
			t.Fatal(cmp.Diff(want, got))
		}
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		defaultSQSEvent,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleEntitlementMessage_SetsProperSEExpressionAttributeValueInUpdateItemWithContextAPICallGivenExpiredEntitlement(t *testing.T) {
	t.Parallel()
	want := &dynamodb.AttributeValue{
		BOOL: aws.Bool(true),
	}
	e, err := entitlement.New("bogus")
	if err != nil {
		t.Fatal(err)
	}
	e.GetEntitlements = func(input *marketplaceentitlementservice.GetEntitlementsInput) (*marketplaceentitlementservice.GetEntitlementsOutput, error) {
		return &marketplaceentitlementservice.GetEntitlementsOutput{
			Entitlements: []*marketplaceentitlementservice.Entitlement{
				{
					ExpirationDate: aws.Time(time.Now().Add(-24 * time.Hour)),
				},
			},
		}, nil
	}
	e.UpdateItemWithContext = func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error) {
		got := input.ExpressionAttributeValues[":se"]
		if !cmp.Equal(want, got) {
			t.Fatal(cmp.Diff(want, got))
		}
		return nil, nil
	}
	err = e.HandleEntitlementMessage(
		context.Background(),
		defaultSQSEvent,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNew_ErrorsGivenEmptyTableName(t *testing.T) {
	t.Parallel()
	_, err := entitlement.New("")
	if err == nil {
		t.Fatal("want error but got nil")
	}
}
