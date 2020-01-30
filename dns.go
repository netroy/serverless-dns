package main

import (
	"encoding/base64"
	"math/rand"
	"net"
	"time"

	"github.com/miekg/dns"
)

// TODO: add a global cache
// https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model-handler-types.html#go-programming-model-handler-execution-environment-reuse

// DNSHandler struct implements #handle to be used as a lambda function
type DNSHandler struct {
	upstreams []string
	client    *dns.Client
}

// NewDNSHandler - Golint made me "document" this
func NewDNSHandler() *DNSHandler {
	handler := new(DNSHandler)

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

func (handler *DNSHandler) randomUpstream() string {
	return handler.upstreams[rand.Intn(len(handler.upstreams))]
}

func (handler *DNSHandler) query(query string) ([]byte, error) {
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
