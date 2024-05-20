package main

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/quicvarint"
	log "github.com/sirupsen/logrus"
)

const (
	link = "https://127.0.0.1:9898/stream"
)

func main() {
	c := NewClient()

	ch1 := make(chan string)
	ch2 := make(chan string)

	go c.CreateStream(link, "Сообщение серверу №1", ch1)
	go c.CreateStream(link, "Сообщение серверу №2", ch2)

	<-ch1
	<-ch2
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

func (c *client) CreateStream(link, message string, ch chan string) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Error(err)
		return
	}
	rtOpt := http3.RoundTripOpt{DontCloseRequestStream: true}
	res, err := c.rt.RoundTripOpt(req, rtOpt)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Client received status %d", res.StatusCode)

	// Открыли стрим
	str := res.Body.(http3.HTTPStreamer).HTTPStream()
	log.Infof("Stream was opened with id=%d", str.StreamID())

	// Отправили данные на сервер
	bytesToSend := []byte(message)
	bytes := quicvarint.Append([]byte{}, uint64(len(bytesToSend)))
	bytes = append(bytes, bytesToSend...)
	if _, err = str.Write(bytes); err != nil {
		log.Error(err)
		return
	}
	log.Infof("Data to server was sent on str id=%d", str.StreamID())

	// Получаем данные от сервера
	lengthOfData, err := quicvarint.Read(quicvarint.NewReader(str))
	bytesToReceive := make([]byte, lengthOfData)
	if _, err = str.Read(bytesToReceive); err != nil {
		log.Error(err)
	}
	log.Infof("Client received data '%s' from server on stream id=%d", string(bytesToReceive), str.StreamID())
	str.CancelRead(quic.StreamErrorCode(quic.NoError))
	ch <- string(bytesToReceive)
}
