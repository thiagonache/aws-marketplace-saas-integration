package landingpage

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/marketplacemetering"
	"github.com/aws/aws-sdk-go/service/sqs"
)

//go:embed templates/*.html templates/*.gohtml
var templates embed.FS

const (
	HTMLSubscribeSuccess = `<div class="alert alert-success" role="alert">
  You have purchased an enterprise product that requires some additional setup.
A representative from our team will be contacting you within two business days with your account credentials.
Please contact Support through our website if you have any questions.
</div>`
)

var (
	ContentTypeTextHTML = map[string]string{
		"content-type": "text/html",
	}
	ContentTypeFormURLEncoded = map[string]string{
		"content-type": "application/x-www-form-urlencoded",
	}
	ErrBadRequest       = errors.New("bad request")
	ErrMethodNotAllowed = errors.New("method not allowed")
	messageBodyTemplate = `{
"Type": "Notification",
"Message" : {
	"action" : "entitlement-updated",
	"customer-identifier": %q,
	"product-code" : %q
	}
}`
	requiredInputs = []string{"inputName", "inputEmail"}
)

type LandingPage struct {
	subscribersTableName       string
	entitlementQueueURL        string
	PutItemWithContext         func(context.Context, *dynamodb.PutItemInput, ...request.Option) (*dynamodb.PutItemOutput, error)
	ResolveCustomerWithContext func(context.Context, *marketplacemetering.ResolveCustomerInput, ...request.Option) (*marketplacemetering.ResolveCustomerOutput, error)
	SendMessageWithContext     func(context.Context, *sqs.SendMessageInput, ...request.Option) (*sqs.SendMessageOutput, error)
}

func New(entitlementQueueURL string, customerTableName string) (LandingPage, error) {
	if entitlementQueueURL == "" || customerTableName == "" {
		return LandingPage{}, errors.New("cannot be empty")
	}
	return LandingPage{
		subscribersTableName: customerTableName,
		entitlementQueueURL:  entitlementQueueURL,
	}, nil
}

func (l LandingPage) CustomerTableName() string {
	return l.subscribersTableName
}

func (l LandingPage) EntitlementQueueURL() string {
	return l.entitlementQueueURL
}

func (l LandingPage) HandleLandingPage(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	switch event.RequestContext.HTTP.Method {
	case http.MethodPost:
		if event.Headers["content-type"] != "application/x-www-form-urlencoded" {
			return events.LambdaFunctionURLResponse{
				StatusCode: http.StatusBadRequest,
				Body:       ErrBadRequest.Error(),
				Headers:    ContentTypeTextHTML,
			}, nil
		}
		token := event.QueryStringParameters["x-amzn-marketplace-token"]
		if token == "" {
			return events.LambdaFunctionURLResponse{
				StatusCode: http.StatusBadRequest,
				Body:       ErrBadRequest.Error(),
				Headers:    ContentTypeTextHTML,
			}, nil
		}
		if event.Body == "" {
			return events.LambdaFunctionURLResponse{
				StatusCode: http.StatusBadRequest,
				Body:       ErrBadRequest.Error(),
				Headers:    ContentTypeTextHTML,
			}, nil
		}
		body := event.Body
		if event.IsBase64Encoded {
			bodyData, err := base64.StdEncoding.DecodeString(event.Body)
			if err != nil {
				return events.LambdaFunctionURLResponse{}, fmt.Errorf("DecodeString: %w", err)
			}
			body = string(bodyData)
		}
		params, err := url.ParseQuery(body)
		if err != nil {
			return events.LambdaFunctionURLResponse{
				StatusCode: http.StatusBadRequest,
				Body:       ErrBadRequest.Error(),
				Headers:    ContentTypeTextHTML,
			}, nil
		}
		for _, input := range requiredInputs {
			if !params.Has(input) {
				return events.LambdaFunctionURLResponse{
					StatusCode: http.StatusBadRequest,
					Body:       ErrBadRequest.Error(),
					Headers:    ContentTypeTextHTML,
				}, nil
			}
		}
		customerOutput, err := l.ResolveCustomerWithContext(ctx, &marketplacemetering.ResolveCustomerInput{
			RegistrationToken: &token,
		})
		if err != nil {
			return events.LambdaFunctionURLResponse{}, err
		}
		companyName := ""
		if len(params["inputCompany"]) > 0 {
			companyName = params["inputCompany"][0]
		}
		contactEmail := params["inputEmail"][0]
		contactJob := ""
		if len(params["inputJob"]) > 0 {
			contactJob = params["inputJob"][0]
		}
		contactName := params["inputName"][0]
		contactPhone := ""
		if len(params["inputPhone"]) > 0 {
			contactPhone = params["inputPhone"][0]
		}
		_, err = l.PutItemWithContext(ctx, &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"companyName":          {S: &companyName},
				"contactEmail":         {S: &contactEmail},
				"contactJob":           {S: &contactJob},
				"contactName":          {S: &contactName},
				"contactPhone":         {S: &contactPhone},
				"lastUpdate":           {S: aws.String(time.Now().UTC().String())},
				"customerAWSAccountID": {S: customerOutput.CustomerAWSAccountId},
				"customerIdentifier":   {S: customerOutput.CustomerIdentifier},
				"productCode":          {S: customerOutput.ProductCode},
			},
			TableName: &l.subscribersTableName,
		})
		if err != nil {
			return events.LambdaFunctionURLResponse{}, err
		}
		msgBody := fmt.Sprintf(messageBodyTemplate, *customerOutput.CustomerIdentifier, *customerOutput.ProductCode)
		_, err = l.SendMessageWithContext(ctx, &sqs.SendMessageInput{
			MessageBody: &msgBody,
			QueueUrl:    &l.entitlementQueueURL,
		})
		if err != nil {
			return events.LambdaFunctionURLResponse{}, err
		}
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusAccepted,
			Body:       HTMLSubscribeSuccess,
			Headers:    ContentTypeTextHTML,
		}, nil
	case http.MethodGet:
		templatePath := "templates/index.gohtml"
		statusCode := http.StatusOK
		marketplaceToken := event.QueryStringParameters["x-amzn-marketplace-token"]
		if marketplaceToken == "" {
			statusCode = http.StatusBadRequest
			templatePath = "templates/badrequest.html"
		}
		tpl, err := template.ParseFS(templates, templatePath)
		if err != nil {
			return events.LambdaFunctionURLResponse{}, err
		}
		htmlPage := &bytes.Buffer{}
		err = tpl.Execute(htmlPage, url.QueryEscape(marketplaceToken))
		if err != nil {
			return events.LambdaFunctionURLResponse{}, err
		}
		return events.LambdaFunctionURLResponse{
			StatusCode: statusCode,
			Body:       htmlPage.String(),
			Headers:    ContentTypeTextHTML,
		}, nil
	default:
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       ErrMethodNotAllowed.Error(),
			Headers:    ContentTypeTextHTML,
		}, nil
	}

}
