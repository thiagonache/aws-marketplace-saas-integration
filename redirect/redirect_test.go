package redirect_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/thiagonache/aws-marketplace-saas-integration/redirect"
)

func TestHandleRedirect_ReturnsStatusFoundGivenPOST(t *testing.T) {
	t.Parallel()
	want := http.StatusFound
	lambdaReq := events.LambdaFunctionURLRequest{
		Headers: map[string]string{
			"content-type": "application/x-www-form-urlencoded",
		},
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/",
			},
		},
		Body:            "eC1hbXpuLW1hcmtldHBsYWNlLXRva2VuPWZha2VUb2tlbg==",
		IsBase64Encoded: true,
	}
	r := redirect.Redirect{}
	lambdaResp, err := r.HandleRedirect(lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.StatusCode
	if want != got {
		t.Fatalf("want response status code %d, got %d", want, got)
	}
}

func TestHandleRedirect_ReturnsProperLocationGivenPOSTWithoutIsBase64EncodedBody(t *testing.T) {
	t.Parallel()
	want := "https://fake.url/landing-page?x-amzn-marketplace-token=fakeToken"
	lambdaReq := events.LambdaFunctionURLRequest{
		Headers: map[string]string{
			"content-type": "application/x-www-form-urlencoded",
		},
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/",
			},
		},
		Body: "x-amzn-marketplace-token=fakeToken",
	}
	u := url.URL{
		Scheme: "https",
		Host:   "fake.url",
		Path:   "/landing-page",
	}
	r := redirect.Redirect{Location: u}
	lambdaResp, err := r.HandleRedirect(lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := lambdaResp.Headers["Location"]
	if !ok {
		t.Fatal("Response Header Location not found")
	}
	if want != got {
		t.Fatalf("want response header Location %q, got %q", want, got)
	}
}

func TestHandleRedirect_ReturnsProperLocationGivenPOSTWithIsBase64EncodedTrueAndValidBody(t *testing.T) {
	t.Parallel()
	want := "https://fake.url/landing-page?x-amzn-marketplace-token=fakeToken"
	lambdaReq := events.LambdaFunctionURLRequest{
		Headers: map[string]string{
			"content-type": "application/x-www-form-urlencoded",
		},
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/",
			},
		},
		Body:            "eC1hbXpuLW1hcmtldHBsYWNlLXRva2VuPWZha2VUb2tlbg==",
		IsBase64Encoded: true,
	}
	u := url.URL{
		Scheme: "https",
		Host:   "fake.url",
		Path:   "/landing-page",
	}
	r := redirect.Redirect{Location: u}
	lambdaResp, err := r.HandleRedirect(lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := lambdaResp.Headers["Location"]
	if !ok {
		t.Fatal("Response Header Location not found")
	}
	if want != got {
		t.Fatalf("want response header Location %q, got %q", want, got)
	}
}

func TestHandleRedirect_ErrorsGivenPOSTWithIsBase64EncodedTrueAndInvalidBody(t *testing.T) {
	t.Parallel()
	lambdaReq := events.LambdaFunctionURLRequest{
		Headers: map[string]string{
			"content-type": "application/x-www-form-urlencoded",
		},
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/",
			},
		},
		Body:            "notBase64EncodedString",
		IsBase64Encoded: true,
	}
	u := url.URL{
		Scheme: "https",
		Host:   "fake.url",
		Path:   "/landing-page",
	}
	r := redirect.Redirect{Location: u}
	_, err := r.HandleRedirect(lambdaReq)
	if err == nil {
		t.Fatal("want error but got nil")
	}
}

func TestHandleRedirect_ReturnsStatusNotAllowedGivenRequestWithUnexpectedMethod(t *testing.T) {
	t.Parallel()
	want := http.StatusMethodNotAllowed
	lambdaReq := events.LambdaFunctionURLRequest{
		Headers: map[string]string{
			"content-type": "application/x-www-form-urlencoded",
		},
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "bogus",
				Path:   "/",
			},
		},
		Body: "body",
	}
	u := url.URL{
		Scheme: "https",
		Host:   "fake.url",
		Path:   "/landing-page",
	}
	r := redirect.Redirect{Location: u}
	lambdaResp, err := r.HandleRedirect(lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.StatusCode
	if want != got {
		t.Fatalf("want response status code %d, got %d", want, got)
	}
}

func TestHandleRedirect_ReturnsStatusBadRequestGivenPOSTWithUnexpectedContentTypeHeader(t *testing.T) {
	t.Parallel()
	want := http.StatusBadRequest
	lambdaReq := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/",
			},
		},
		Headers: map[string]string{
			"content-type": "bogus",
		},
	}
	u := url.URL{
		Scheme: "https",
		Host:   "fake.url",
		Path:   "/landing-page",
	}
	r := redirect.Redirect{Location: u}
	lambdaResp, err := r.HandleRedirect(lambdaReq)
	if err != nil {
		t.Fatal(err)
	}
	got := lambdaResp.StatusCode
	if want != got {
		t.Fatalf("want response status code %d, got %d", want, got)
	}
}
