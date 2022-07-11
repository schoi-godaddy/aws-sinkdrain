package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Headers        map[string]interface{}
	RequestContext map[string]interface{}
	Body           map[string]interface{}
}

type APIGatewayEvent struct {
	Headers         map[string]interface{} `json:"headers"`
	RequestContext  map[string]interface{} `json:"requestContext"`
	Body            string                 `json:"body"`
	IsBase64Encoded bool                   `json:"isBase64Encoded"`
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// https://www.crimsonmacaw.com/blog/handling-multiple-aws-lambda-event-types-with-go/#aws-lambda-in-go
// https://pkg.go.dev/encoding/json#Unmarshaler
func (event *Event) UnmarshalJSON(data []byte) error {
	// log.Print("data ", string(data))
	e := &APIGatewayEvent{}

	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}

	body := &map[string]interface{}{}

	if e.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(e.Body)

		if err != nil {
			return err
		}

		if err := json.Unmarshal(b, &body); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal([]byte(e.Body), &body); err != nil {
			return err
		}
	}

	event.Body = *body
	event.Headers = e.Headers
	event.RequestContext = e.RequestContext

	return nil
}

func HandleRequest(ctx context.Context, e Event) (Response, error) {
	divisor := 2
	MaskMap(e.Body, divisor)

	body, err := json.Marshal(e.Body)

	if err != nil {
		log.Print("Error trying to marshal e")
		return Response{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "*",
				"Access-Control-Allow-Headers": "*",
			},
			Body: "Error trying to marshal event",
		}, err
	}

	log.Print("Body - ", string(body))
	// log.Print("Headers - ", e.Headers)
	// log.Print("RequestContext - ", e.RequestContext)

	return Response{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
		},
		Body: "Function success!",
	}, nil
}

// Map are pass by reference in Go
func MaskMap(m map[string]interface{}, divisor int) {
	for key, val := range m {
		if v, ok := val.(string); ok {
			m[key] = fmt.Sprintf("%s%s", v[:len(v)/divisor], strings.Repeat("-", len(v)-len(v)/divisor))
		}
	}
}

func main() {
	lambda.Start(HandleRequest)
}
