package packages

import (
	"crypto/rand"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/ory/graceful"

	"github.com/restuwahyu13/discovery-api/helpers"
)

type GracefulConfig struct {
	Handler *chi.Mux
	Port    string
}

func Graceful(Handler func() *GracefulConfig) error {
	parser := helpers.NewParser()
	inboundSize, _ := parser.ToInt(os.Getenv("INBOUND_SIZE"))

	h := Handler()
	secure := true

	if _, ok := os.LookupEnv("GO_ENV"); ok && os.Getenv("GO_ENV") != "development" {
		secure = false
	}

	server := http.Server{
		Handler:        h.Handler,
		Addr:           ":" + h.Port,
		MaxHeaderBytes: inboundSize,
		TLSConfig: &tls.Config{
			Rand:               rand.Reader,
			InsecureSkipVerify: secure,
		},
	}

	Logrus("info", "Server listening on port %s", h.Port)
	return graceful.Graceful(server.ListenAndServe, server.Shutdown)
}
