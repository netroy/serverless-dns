package main

import (
	"encoding/base64"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestDNSHandler(t *testing.T) {
	dnsHandler := NewDNSHandler(&Config{
		timeout:   1,
		upstreams: []string{"94.140.15.15:53"},
	})
	request := new(dns.Msg)
	request.SetQuestion(dns.Fqdn("one.one.one.one"), dns.TypeA)
	binary, _ := request.Pack()

	query := base64.RawURLEncoding.EncodeToString(binary)
	answer, _ := dnsHandler.Query(query)
	response := new(dns.Msg)
	response.Unpack(answer)

	answers := make([]string, len(response.Answer))
	for i, v := range response.Answer {
		answers[i] = v.(*dns.A).A.String()
	}

	expected := []string{"1.1.1.1", "1.0.0.1"}
	assert.ElementsMatch(t, answers, expected)
}
