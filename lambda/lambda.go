package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func parseQuery(req events.APIGatewayProxyRequest) (string, error) {
	method := req.HTTPMethod
	switch method {
	case "GET":
		return req.QueryStringParameters["dns"], nil
	case "POST":
		return req.Body, nil
	}
	return "", errors.New("Invalid DNS request")
}

// LambdaHandler returns function that can be used by AWS Lambda runtime
func LambdaHandler(dnsHandler *DNSHandler) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		outputJSON, _ := json.Marshal(req)
		log.Println(outputJSON)

		query, err := parseQuery(req)
		if err != nil || len(query) == 0 {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, nil
		}

		body, err := dnsHandler.Query(query)
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
	lambdaHandler := LambdaHandler(dnsHandler)
	lambda.Start(lambdaHandler)
	// query := "JGUBAAABAAAAAAAAB2V4YW1wbGUDY29tAAABAAE"
	// params := map[string]string{query: query}
	// req := &events.APIGatewayProxyRequest{
	// 	QueryStringParameters: params,
	// }
	// resp, err := lambdaHandler(nil, *req)
	// log.Println(resp, err)
}
