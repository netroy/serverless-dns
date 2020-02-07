package main

import (
	"encoding/base64"
	"math/rand"
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/netroy/ttlcache"
)

// TODO: add a global cache
// https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model-handler-types.html#go-programming-model-handler-execution-environment-reuse

// CacheEntry wraps DNS responses in a struct, so we can cache them
type CacheEntry struct {
	data []byte
}

// DNSHandler struct implements #handle to be used as a lambda function
type DNSHandler struct {
	upstreams []string
	client    *dns.Client
	cache     *ttlcache.Cache
}

// NewDNSHandler - Golint made me "document" this
func NewDNSHandler() *DNSHandler {
	handler := new(DNSHandler)

	// TODO: don't hard code
	handler.cache = ttlcache.NewCache(60 * time.Second)

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

// TODO: allow chosing between random & round-robin
func (handler *DNSHandler) randomUpstream() string {
	return handler.upstreams[rand.Intn(len(handler.upstreams))]
}

// Query lets you resolve DNS queries.
// This uses a TTL cache to reduce querying upstream
func (handler *DNSHandler) Query(query string) ([]byte, error) {
	entry, exists := handler.cache.Get(query)
	if exists == true {
		return entry.(CacheEntry).data, nil
	}

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

	data, err := response.Pack()
	if err != nil {
		return nil, err
	}

	handler.cache.Set(query, &CacheEntry{data})
	return binary, nil
}
