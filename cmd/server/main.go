package main

import (
	"net/http"

	"github.com/quic-go/quic-go/http3"
	log "github.com/sirupsen/logrus"
)

const (
	addr = "127.0.0.1:9898"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/stream", HandleStream)
	srv := http3.Server{
		Addr: addr,
	}
	if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		log.Error(err)
	}
}

func HandleStream(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()

	// Обработка стрима
}