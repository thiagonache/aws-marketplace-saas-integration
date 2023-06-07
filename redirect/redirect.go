package redirect

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

var (
	ContentTypeTextHTML = map[string]string{
		"content-type": "text/html",
	}
	ErrBadRequest       = errors.New("bad request")
	ErrMethodNotAllowed = errors.New("method not allowed")
)

type Redirect struct {
	Location url.URL
}

func (r Redirect) HandleRedirect(event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	switch event.RequestContext.HTTP.Method {
	case http.MethodPost:
		if event.Headers["content-type"] != "application/x-www-form-urlencoded" {
			return events.LambdaFunctionURLResponse{
				StatusCode: http.StatusBadRequest,
				Body:       ErrBadRequest.Error(),
				Headers:    ContentTypeTextHTML,
			}, nil
		}
		marketplaceToken := event.Body
		if event.IsBase64Encoded {
			tokenData, err := base64.StdEncoding.DecodeString(event.Body)
			if err != nil {
				return events.LambdaFunctionURLResponse{}, err
			}
			marketplaceToken = string(tokenData)
		}
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusFound,
			Headers: map[string]string{
				"Location": r.Location.String() + "?" + marketplaceToken,
			},
		}, nil
	default:
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       ErrMethodNotAllowed.Error(),
		}, nil
	}
}
