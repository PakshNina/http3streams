package main

import (
	"fmt"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/quicvarint"
	log "github.com/sirupsen/logrus"
)

const (
	addr = "127.0.0.1:9898"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/stream", HandleStream)
	srv := http3.Server{
		Addr:    addr,
		Handler: mux,
	}
	if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		log.Error(err)
	}
}

func HandleStream(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()

	// Обработка стрима
	str := req.Body.(http3.HTTPStreamer).HTTPStream()
	log.Infof("Stream with id=%d from client was accepted", str.StreamID())
	go readAndWrite(str)
}

func readAndWrite(str quic.Stream) {
	// Считываем данные от клиента
	lengthOfData, err := quicvarint.Read(quicvarint.NewReader(str))
	if err != nil {
		log.Error(err)
		return
	}
	bytes := make([]byte, lengthOfData)
	if _, err = str.Read(bytes); err != nil {
		log.Error(err)
		return
	}
	log.Infof("Client sent data '%s' from stream id=%d", string(bytes), str.StreamID())

	// Отправляем данные клиенту
	message := []byte(fmt.Sprintf("Ответ сервера на стриме id=%d", str.StreamID()))
	bytesToSend := quicvarint.Append([]byte{}, uint64(len(message)))
	bytesToSend = append(bytesToSend, message...)

	if _, err = str.Write(bytesToSend); err != nil {
		log.Error(err)
		return
	}
	<-str.Context().Done()
}
