package app

import (
	"log"
	"net/http"
)

func (a *App) Run(addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: a.Router,
	}

	log.Printf("server listening on http://localhost%s", addr)

	return server.ListenAndServe()
}
