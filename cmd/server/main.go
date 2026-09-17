package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/sasvyn/backend/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	port := 8080
	if configuredPort := os.Getenv("PORT"); configuredPort != "" {
		parsedPort, err := strconv.Atoi(configuredPort)
		if err != nil || parsedPort < 1 || parsedPort > 65535 {
			log.Fatalf("invalid PORT %q", configuredPort)
		}
		port = parsedPort
	}

	server := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           newHandler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	go func() {
		<-shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("server shutdown: %v", err)
		}
	}()

	log.Printf("server listening on http://localhost:%d", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server: %v", err)
	}
}

func newHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /vijay", vijayHandler)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}

func healthHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(map[string]string{"status": "ok"})
}

// vijayHandler godoc
// @Summary Test Vijay endpoint
// @Description Returns a simple test response
// @Tags Test
// @Produce json
// @Success 200 {object} map[string]string
// @Router /vijay [get]
func vijayHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(map[string]string{
		"status": "ok",
		"data":   "vijay",
	})
}