package landingpage_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/marketplacemetering"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/go-cmp/cmp"
	"github.com/thiagonache/aws-marketplace-saas-integration/landingpage"
)

func TestHandleLandingPage_RendersProperHTMLGivenGETWithMarketplaceToken(t *testing.T) {
	t.Parallel()
	want := `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Marketplace Registration Page</title>
    <link
      href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css"
      rel="stylesheet"
      integrity="sha384-KK94CHFLLe+nY2dmCWGMq91rCGa5gtU4mk92HdvYe+M/SXH301p5ILy+dN9+nJOZ"
      crossorigin="anonymous"
    />
  </head>
  <body class="text-center">
    <div class="container-lg">
      <nav class="navbar bg-dark">
        <div class="container-fluid">
          <h1>Your Company Name</h1>
        </div>
      </nav>
      <h1>Registration Information</h1>
      <form hx-post="/?x-amzn-marketplace-token=bogus">
        <div class="input-group mb-3">
          <span class="input-group-text">Company</span>
          <input
            type="text"
            class="form-control"
            name="inputCompany"
            aria-describedby="companyHelp"
            autofocus
          />
        </div>
        <div class="input-group mb-3">
          <span class="input-group-text">Name</span>
          <input
            type="text"
            class="form-control"
            name="inputName"
            aria-describedby="nameHelp"
            required
          />
        </div>
        <div class="input-group mb-3">
          <span class="input-group-text">Job</span>
          <input
            type="text"
            class="form-control"
            name="inputJob"
            aria-describedby="jobHelp"
          />
        </div>
        <div class="input-group mb-3">
          <span class="input-group-text">Email</span>
          <input
            type="email"
            class="form-control"
            name="inputEmail"
            aria-describedby="emailHelp"
            required
          />
        </div>
        <div class="input-group mb-3">
          <span class="input-group-text">Phone</span>
          <input type="text" name="inputPhone" class="form-control" />
        </div>
        <div class="mb-3">
          <button class="btn btn-dark">Submit</button>
        </div>
      </form>
      <p class="mt-5 mb-3 text-muted">
        &copy; 2023 Thiago Nache Carvalho, Inc. All Rights Reserved
      </p>
    </div>
    <script src="https://unpkg.com/htmx.org@1.9.2"></script>
  </body>
</html>
`
	lambdaReq := events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "bogus",
		},
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodGet,
				Path:   "/",
			},
		},
	}
	l := landingpage.LandingPage{}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.Body
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestHandleLandingPage_ReturnsStatusOKGivenGETWithMarketplaceToken(t *testing.T) {
	t.Parallel()
	want := http.StatusOK
	lambdaReq := events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "bogus",
		},
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodGet,
				Path:   "/",
			},
		},
	}
	l := landingpage.LandingPage{}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.StatusCode
	if want != got {
		t.Fatalf("want response status code %d, got %d", want, got)
	}
}

func TestHandleLandingPage_ReturnsContentTypeTextHTMLGivenGETWithMarketplaceToken(t *testing.T) {
	t.Parallel()
	wantHeader := "content-type"
	want := "text/html"
	lambdaReq := events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "bogus",
		},
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodGet,
				Path:   "/",
			},
		},
	}
	l := landingpage.LandingPage{}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.Headers[wantHeader]
	if want != got {
		t.Fatalf("want header %q %q, got %q", wantHeader, want, got)
	}
}

func TestHandleLandingPage_RendersProperHTMLGivenGETWithoutMarketplaceToken(t *testing.T) {
	t.Parallel()
	want := `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Marketplace Registration Page</title>
    <link
      href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css"
      rel="stylesheet"
      integrity="sha384-KK94CHFLLe+nY2dmCWGMq91rCGa5gtU4mk92HdvYe+M/SXH301p5ILy+dN9+nJOZ"
      crossorigin="anonymous"
    />
  </head>
  <body class="text-center">
    <div class="container-lg">
      <nav class="navbar bg-dark">
        <div class="container-fluid">
          <h1>Your Company Name</h1>
        </div>
      </nav>
      <h1>Registration Information</h1>
      <div class="alert alert-danger" role="alert">
        Please, go to AWS Marketplace to Set up your product!
      </div>
      <p class="mt-5 mb-3 text-muted text-center">
        &copy; 2023 Thiago Nache Carvalho, Inc. All Rights Reserved
      </p>
    </div>
  </body>
</html>
`
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodGet,
				Path:   "/",
			},
		},
	}
	l := landingpage.LandingPage{}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.Body
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestHandleLandingPage_ReturnsStatusBadRequestGivenGETWithoutMarketplaceToken(t *testing.T) {
	t.Parallel()
	want := http.StatusBadRequest
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodGet,
				Path:   "/",
			},
		},
		Headers: landingpage.ContentTypeFormURLEncoded,
	}
	l := landingpage.LandingPage{}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.StatusCode
	if want != got {
		t.Fatalf("want response status code %d, got %d", want, got)
	}
}

func TestHandleLandingPage_ReturnsContentTypeTextHTMLGivenGETWithoutMarketplaceToken(t *testing.T) {
	t.Parallel()
	wantHeader := "content-type"
	want := "text/html"
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodGet,
				Path:   "/",
			},
		},
		Headers: landingpage.ContentTypeFormURLEncoded,
	}
	l := landingpage.LandingPage{}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.Headers[wantHeader]
	if want != got {
		t.Fatalf("want header %q %q, got %q", wantHeader, want, got)
	}
}

func TestHandleLandingPage_ReturnsStatusAcceptedGivenPOSTWithBody(t *testing.T) {
	t.Parallel()
	want := http.StatusAccepted
	body := base64.StdEncoding.EncodeToString([]byte("inputName=bogus&inputEmail=bogus"))
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
				Path:   "/",
			},
		},
		Headers: landingpage.ContentTypeFormURLEncoded,
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "bogus",
		},
		Body:            body,
		IsBase64Encoded: true,
	}
	l := landingpage.LandingPage{
		ResolveCustomerWithContext: func(ctx context.Context, input *marketplacemetering.ResolveCustomerInput, opts ...request.Option) (*marketplacemetering.ResolveCustomerOutput, error) {
			return &marketplacemetering.ResolveCustomerOutput{
				CustomerIdentifier: new(string),
				ProductCode:        new(string),
			}, nil
		},
		SendMessageWithContext: func(ctx context.Context, input *sqs.SendMessageInput, opts ...request.Option) (*sqs.SendMessageOutput, error) {
			return nil, nil
		},
		PutItemWithContext: func(ctx context.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
			return nil, nil
		},
	}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.StatusCode
	if want != got {
		t.Fatalf("want response status code %d, got %d", want, got)
	}
}

func TestHandleLandingPage_ReturnsStatusBadRequestGivenPOSTWithUnexpectedContentTypeHeader(t *testing.T) {
	t.Parallel()
	want := http.StatusBadRequest
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
				Path:   "/",
			},
		},
		Headers: map[string]string{
			"content-type": "bogus",
		},
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "bogus",
		},
	}
	l := landingpage.LandingPage{}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.StatusCode
	if want != got {
		t.Fatalf("want response status code %d, got %d", want, got)
	}
}

func TestHandleLandingPage_ReturnsStatusBadRequestGivenPOSTMissing(t *testing.T) {
	testCases := []struct {
		desc string
		body string
	}{
		{
			desc: "everything in body",
			body: "",
		},
		{
			desc: "inputName in body",
			body: "inputEmail=bogus",
		},
		{
			desc: "inputEmail in body",
			body: "inputName=bogus",
		},
	}
	want := http.StatusBadRequest
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			lambdaReq := events.LambdaFunctionURLRequest{
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
						Method: http.MethodPost,
						Path:   "/",
					},
				},
				Headers: landingpage.ContentTypeFormURLEncoded,
				Body:    tC.body,
				QueryStringParameters: map[string]string{
					"x-amzn-marketplace-token": "bogus",
				},
			}
			l := landingpage.LandingPage{}
			lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
			if err != nil {
				t.Fatal(err)
			}
			got := lambdaResp.StatusCode
			if want != got {
				t.Fatalf("want response status code %d, got %d", want, got)
			}
		})
	}
}

func TestHandleLandingPage_ReturnsStatusBadRequestGivenPOSTWithoutMarketplaceToken(t *testing.T) {
	t.Parallel()
	want := http.StatusBadRequest
	body := base64.StdEncoding.EncodeToString([]byte("inputName=bogus&inputEmail=bogus"))
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
				Path:   "/",
			},
		},
		Headers:         landingpage.ContentTypeFormURLEncoded,
		Body:            body,
		IsBase64Encoded: true,
	}
	l := landingpage.LandingPage{}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.StatusCode
	if want != got {
		t.Fatalf("want response status code %d, got %d", want, got)
	}
}

func TestHandleLandingPage_ReturnsStatusMethodNotAllowedGivenRequestWithUnexpectedMethod(t *testing.T) {
	t.Parallel()
	want := http.StatusMethodNotAllowed
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "bogus",
				Path:   "/",
			},
		},
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "bogus",
		},
	}
	l := landingpage.LandingPage{}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.StatusCode
	if want != got {
		t.Fatalf("want response status code %d, got %d", want, got)
	}
}

func TestHandleLandingPage_CallsResolveCustomerWithContext(t *testing.T) {
	t.Parallel()
	called := false
	body := base64.StdEncoding.EncodeToString([]byte("inputName=bogus&inputEmail=bogus"))
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
				Path:   "/",
			},
		},
		Headers: landingpage.ContentTypeFormURLEncoded,
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "bogus",
		},
		Body:            body,
		IsBase64Encoded: true,
	}
	l := landingpage.LandingPage{
		ResolveCustomerWithContext: func(ctx context.Context, input *marketplacemetering.ResolveCustomerInput, opts ...request.Option) (*marketplacemetering.ResolveCustomerOutput, error) {
			called = true
			return &marketplacemetering.ResolveCustomerOutput{
				CustomerIdentifier: new(string),
				ProductCode:        new(string),
			}, nil
		},
		SendMessageWithContext: func(ctx context.Context, smi *sqs.SendMessageInput, o ...request.Option) (*sqs.SendMessageOutput, error) {
			return nil, nil
		},
		PutItemWithContext: func(ctx context.Context, pii *dynamodb.PutItemInput, o ...request.Option) (*dynamodb.PutItemOutput, error) {
			return nil, nil
		},
	}
	_, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("function ResolveCustomerWithContext not called")
	}
}

func TestHandleLandingPage_SetsProperRegistrationTokenInResolveCustomerWithContextAPICall(t *testing.T) {
	t.Parallel()
	want := "CallsResolveCustomerWithCorrectRegistrationToken"
	body := base64.StdEncoding.EncodeToString([]byte("inputName=bogus&inputEmail=bogus"))
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
				Path:   "/",
			},
		},
		Headers: landingpage.ContentTypeFormURLEncoded,
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "CallsResolveCustomerWithCorrectRegistrationToken",
		},
		Body:            body,
		IsBase64Encoded: true,
	}
	l := landingpage.LandingPage{
		ResolveCustomerWithContext: func(ctx context.Context, input *marketplacemetering.ResolveCustomerInput, opts ...request.Option) (*marketplacemetering.ResolveCustomerOutput, error) {
			got := *input.RegistrationToken
			if want != got {
				t.Fatalf("want registration token %q, got %q", want, got)
			}
			return &marketplacemetering.ResolveCustomerOutput{
				CustomerIdentifier: new(string),
				ProductCode:        new(string),
			}, nil
		},
		SendMessageWithContext: func(ctx context.Context, smi *sqs.SendMessageInput, o ...request.Option) (*sqs.SendMessageOutput, error) {
			return nil, nil
		},
		PutItemWithContext: func(ctx context.Context, pii *dynamodb.PutItemInput, o ...request.Option) (*dynamodb.PutItemOutput, error) {
			return nil, nil
		},
	}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	if lambdaResp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected response status code %d", lambdaResp.StatusCode)
	}
}

func TestHandleLandingPage_CallsSendMessageWithContext(t *testing.T) {
	t.Parallel()
	called := false
	body := base64.StdEncoding.EncodeToString([]byte("inputName=bogus&inputEmail=bogus"))
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
				Path:   "/",
			},
		},
		Headers: landingpage.ContentTypeFormURLEncoded,
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "bogus",
		},
		Body:            body,
		IsBase64Encoded: true,
	}
	l := landingpage.LandingPage{
		ResolveCustomerWithContext: func(ctx context.Context, input *marketplacemetering.ResolveCustomerInput, opts ...request.Option) (*marketplacemetering.ResolveCustomerOutput, error) {
			return &marketplacemetering.ResolveCustomerOutput{
				CustomerIdentifier: new(string),
				ProductCode:        new(string),
			}, nil
		},
		SendMessageWithContext: func(ctx context.Context, smi *sqs.SendMessageInput, o ...request.Option) (*sqs.SendMessageOutput, error) {
			called = true
			return nil, nil
		},
		PutItemWithContext: func(ctx context.Context, pii *dynamodb.PutItemInput, o ...request.Option) (*dynamodb.PutItemOutput, error) {
			return nil, nil
		},
	}
	_, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("function ResolveCustomerWithContext not called")
	}
}

func TestHandleLandingPage_SetsProperMessageBodyInSendMessageWithContextAPICall(t *testing.T) {
	t.Parallel()
	want := `{
"Type": "Notification",
"Message" : {
	"action" : "entitlement-updated",
	"customer-identifier": "anyGlobalUniqueIdentifierGivenByAWS",
	"product-code" : "mySAASProductCodeGivenByAWS"
	}
}`
	body := base64.StdEncoding.EncodeToString([]byte("inputName=bogus&inputEmail=bogus"))
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
				Path:   "/",
			},
		},
		Headers: landingpage.ContentTypeFormURLEncoded,
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "bogus",
		},
		Body:            body,
		IsBase64Encoded: true,
	}
	l := landingpage.LandingPage{
		ResolveCustomerWithContext: func(ctx context.Context, input *marketplacemetering.ResolveCustomerInput, opts ...request.Option) (*marketplacemetering.ResolveCustomerOutput, error) {
			return &marketplacemetering.ResolveCustomerOutput{
				CustomerIdentifier: aws.String("anyGlobalUniqueIdentifierGivenByAWS"),
				ProductCode:        aws.String("mySAASProductCodeGivenByAWS"),
			}, nil
		},
		SendMessageWithContext: func(ctx context.Context, input *sqs.SendMessageInput, opts ...request.Option) (*sqs.SendMessageOutput, error) {
			got := *input.MessageBody
			if !cmp.Equal(want, got) {
				t.Fatal(cmp.Diff(want, got))
			}
			return nil, nil
		},
		PutItemWithContext: func(ctx context.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
			return nil, nil
		},
	}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	if lambdaResp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected response status code %d", lambdaResp.StatusCode)
	}
}

func TestHandleLandingPage_SetsProperEntitlementQueueURLInSendMessageWithContextAPICall(t *testing.T) {
	t.Parallel()
	want := "https://sqs.us-east-1.amazonaws.com/177715257436/MyQueue"
	body := base64.StdEncoding.EncodeToString([]byte("inputName=bogus&inputEmail=bogus"))
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
				Path:   "/",
			},
		},
		Headers: landingpage.ContentTypeFormURLEncoded,
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "bogus",
		},
		Body:            body,
		IsBase64Encoded: true,
	}
	l, err := landingpage.New("https://sqs.us-east-1.amazonaws.com/177715257436/MyQueue", "bogus")
	if err != nil {
		t.Fatal(err)
	}
	l.ResolveCustomerWithContext = func(ctx context.Context, input *marketplacemetering.ResolveCustomerInput, opts ...request.Option) (*marketplacemetering.ResolveCustomerOutput, error) {
		return &marketplacemetering.ResolveCustomerOutput{
			CustomerIdentifier: aws.String("anyGlobalUniqueIdentifierGivenByAWS"),
			ProductCode:        aws.String("mySAASProductCodeGivenByAWS"),
		}, nil
	}
	l.SendMessageWithContext = func(ctx context.Context, input *sqs.SendMessageInput, opts ...request.Option) (*sqs.SendMessageOutput, error) {
		got := *input.QueueUrl
		if want != got {
			t.Fatalf("want queue url %q, got %q", want, got)
		}
		return nil, nil
	}
	l.PutItemWithContext = func(ctx context.Context, pii *dynamodb.PutItemInput, o ...request.Option) (*dynamodb.PutItemOutput, error) {
		return nil, nil
	}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	if lambdaResp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected response status code %d", lambdaResp.StatusCode)
	}
}

func TestHandleLandingPage_CallsPutItemWithContext(t *testing.T) {
	t.Parallel()
	called := false
	body := base64.StdEncoding.EncodeToString([]byte("inputName=bogus&inputEmail=bogus"))
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
				Path:   "/",
			},
		},
		Headers: landingpage.ContentTypeFormURLEncoded,
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "bogus",
		},
		Body:            body,
		IsBase64Encoded: true,
	}
	l := landingpage.LandingPage{
		ResolveCustomerWithContext: func(ctx context.Context, input *marketplacemetering.ResolveCustomerInput, opts ...request.Option) (*marketplacemetering.ResolveCustomerOutput, error) {
			return &marketplacemetering.ResolveCustomerOutput{
				CustomerIdentifier: new(string),
				ProductCode:        new(string),
			}, nil
		},
		SendMessageWithContext: func(ctx context.Context, smi *sqs.SendMessageInput, o ...request.Option) (*sqs.SendMessageOutput, error) {
			return nil, nil
		},
		PutItemWithContext: func(ctx context.Context, pii *dynamodb.PutItemInput, o ...request.Option) (*dynamodb.PutItemOutput, error) {
			called = true
			return nil, nil
		},
	}
	_, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("function ResolveCustomerWithContext not called")
	}
}

func TestHandleLandingPage_SetsProperTableNameInPutItemWithContextAPICall(t *testing.T) {
	t.Parallel()
	want := "MyDynamoDBTable"
	body := base64.StdEncoding.EncodeToString([]byte("inputName=bogus&inputEmail=bogus"))
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
				Path:   "/",
			},
		},
		Headers: landingpage.ContentTypeFormURLEncoded,
		QueryStringParameters: map[string]string{
			"x-amzn-marketplace-token": "bogus",
		},
		Body:            body,
		IsBase64Encoded: true,
	}
	l, err := landingpage.New("bogus", "MyDynamoDBTable")
	if err != nil {
		t.Fatal(err)
	}
	l.ResolveCustomerWithContext = func(ctx context.Context, input *marketplacemetering.ResolveCustomerInput, opts ...request.Option) (*marketplacemetering.ResolveCustomerOutput, error) {
		return &marketplacemetering.ResolveCustomerOutput{
			CustomerIdentifier: new(string),
			ProductCode:        new(string),
		}, nil
	}
	l.SendMessageWithContext = func(ctx context.Context, input *sqs.SendMessageInput, opts ...request.Option) (*sqs.SendMessageOutput, error) {
		return nil, nil
	}
	l.PutItemWithContext = func(ctx context.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
		got := *input.TableName
		if want != got {
			t.Fatalf("want table name %q, got %q", want, got)
		}
		return nil, nil
	}
	lambdaResp, err := l.HandleLandingPage(context.Background(), lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	if lambdaResp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected response status code %d", lambdaResp.StatusCode)
	}
}

func TestNew_SetsExpectedEntitlementQueueURL(t *testing.T) {
	want := "https://sqs.us-east-1.amazonaws.com/177715257436/MyQueue"
	l, err := landingpage.New("https://sqs.us-east-1.amazonaws.com/177715257436/MyQueue", "bogus")
	if err != nil {
		t.Fatal(err)
	}
	got := l.EntitlementQueueURL()
	if want != got {
		t.Fatalf("want queue URL %q, got %q", want, got)
	}
}

func TestNew_ErrorsGivenEmptyEntitlementQueueURL(t *testing.T) {
	t.Parallel()
	_, err := landingpage.New("", "bogus")
	if err == nil {
		t.Fatal("want error but got nil")
	}
}

func TestNew_SetsExpectedEntitlementDynamoDBTableName(t *testing.T) {
	want := "MyDynamoDBTable"
	l, err := landingpage.New("bogus", "MyDynamoDBTable")
	if err != nil {
		t.Fatal(err)
	}
	got := l.CustomerTableName()
	if want != got {
		t.Fatalf("want queue URL %q, got %q", want, got)
	}
}

func TestNew_ErrorsGivenEmptyDynamoDBTableName(t *testing.T) {
	t.Parallel()
	_, err := landingpage.New("bogus", "")
	if err == nil {
		t.Fatal("want error but got nil")
	}
}
