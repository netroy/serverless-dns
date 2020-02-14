package main

import (
	"encoding/base64"
	"log"
	"testing"

	"github.com/miekg/dns"
)

// m := new(dns.Msg)
// m.SetQuestion(dns.Fqdn("example.com."), dns.TypeA)
// binary, _ := m.Pack()
// query := base64.RawURLEncoding.EncodeToString(binary)

func TestQuery(t *testing.T) {
	query := "JGUBAAABAAAAAAAAB2V4YW1wbGUDY29tAAABAAE"
	binary, _ := base64.RawURLEncoding.DecodeString(query)
	msg := new(dns.Msg)
	msg.Unpack(binary)

	log.Println(msg.Question[0].Name)
}
