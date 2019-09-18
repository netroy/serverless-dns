package main

import (
	"context"
	"encoding/base64"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/miekg/dns"
)

// Handler struct implements #handle to be used as a lambda function
type Handler struct {
	upstreams []string
	client    *dns.Client
}

// NewHandler - Golint made me "document" this
func NewHandler() *Handler {
	handler := new(Handler)

	resolv, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	count := len(resolv.Servers)
	if count > 0 {
		handler.upstreams = make([]string, count)
		for index, server := range resolv.Servers {
			handler.upstreams[index] = server + ":53"
		}
	} else {
		// TODO: don't hardcode
		handler.upstreams = []string{"1.1.1.1"}
	}

	// TODO: don't hardcode
	timeout := time.Duration(15) * time.Second
	handler.client = &dns.Client{
		Net:     "tcp",
		Timeout: timeout,
	}
	handler.client.Dialer = &net.Dialer{
		Timeout:   timeout,
		LocalAddr: nil,
	}

	return handler
}

func (handler *Handler) randomUpstream() string {
	return handler.upstreams[rand.Intn(len(handler.upstreams))]
}

func (handler *Handler) query(query string) ([]byte, error) {
	binary, err := base64.RawURLEncoding.DecodeString(query)
	if err != nil {
		return nil, err
	}
	msg := new(dns.Msg)
	msg.Unpack(binary)

	upstream := handler.randomUpstream()
	response, _, err := handler.client.Exchange(msg, upstream)
	if err != nil {
		return nil, err
	}

	binary, err = response.Pack()
	if err != nil {
		return nil, err
	}
	return binary, nil
}

func (handler *Handler) handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dns := req.QueryStringParameters["dns"]
	body, err := handler.query(dns)
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

func main() {
	handler := NewHandler()
	lambda.Start(handler.handle)
}
