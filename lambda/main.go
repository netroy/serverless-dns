package main

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/netroy/ttlcache"
)

func LambdaHandler(dnsHandler *DNSHandler) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

var dnsHandler *DNSHandler

func init() {
	if dnsHandler == nil {
		dnsHandler = NewDNSHandler(&Config{
			cache:   ttlcache.NewCache(60 * time.Second),
			timeout: 3,
			upstreams: []string{
				// AdGuard
				"94.140.14.14",
				"94.140.15.15",

				// Cloudflare
				// "1.1.1.1:53",
				// "1.0.0.1:53",

				// Google
				// "8.8.8.8:53",
				// "8.8.4.4:53",
			},
		})
	}
}

func main() {
	lambda.Start(LambdaHandler(dnsHandler))
}
