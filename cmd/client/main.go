package main

import (
	"crypto/tls"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

const (
	link = "https://127.0.0.1:9898/stream"
)

func main() {

}

type client struct {
	rt *http3.RoundTripper
}

func NewClient() *client {
	rt := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		QuicConfig: &quic.Config{
			MaxIdleTimeout: time.Minute,
		},
	}
	return &client{rt: rt}
}
