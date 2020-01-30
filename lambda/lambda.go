package main

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func LambdaHandler(dnsHandler *DNSHandler) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		dns := req.QueryStringParameters["dns"]
		body, err := dnsHandler.query(dns)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       base64.StdEncoding.EncodeToString(body),
			Headers: map[string]string{
				"Content-Type":   "application/dns-message",
				"Content-Length": strconv.Itoa(len(body)),
				// TODO: handle TTL as Expires
			},
			IsBase64Encoded: true,
		}, nil
	}
}

func main() {
	dnsHandler := NewDNSHandler()
	lambda.Start(LambdaHandler(dnsHandler))
}
