package main

import (
	"encoding/base64"
	"math/rand"
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/netroy/ttlcache"
)

// CacheEntry wraps DNS responses in a struct, so we can cache them
type CacheEntry struct {
	data []byte
}

type Config struct {
	cache     *ttlcache.Cache
	timeout   int
	upstreams []string
}

type DNSHandler struct {
	upstreams []string
	client    *dns.Client
	cache     *ttlcache.Cache
}

// NewDNSHandler - Golint made me "document" this
func NewDNSHandler(config *Config) *DNSHandler {
	handler := new(DNSHandler)
	handler.cache = config.cache
	handler.upstreams = config.upstreams
	timeout := time.Duration(config.timeout) * time.Second
	handler.client = &dns.Client{
		Net:     "udp",
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
	if handler.cache != nil {
		entry, exists := handler.cache.Get(query)
		if exists == true {
			return entry.(*CacheEntry).data, nil
		}
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
	if handler.cache != nil {
		handler.cache.Set(query, &CacheEntry{data})
	}
	return data, nil
}
