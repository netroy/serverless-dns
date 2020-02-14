package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// LambdaHandler returns function that can be used by AWS Lambda runtime
func LambdaHandler(dnsHandler *DNSHandler) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		dns := req.QueryStringParameters["dns"]
		body, err := dnsHandler.Query(dns)
		if err != nil {
			outputJSON, _ := json.Marshal(req)
			log.Println(outputJSON)
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
