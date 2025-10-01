package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"ntp-service/handlers"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	http.HandleFunc("/time", handlers.TimeHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	serverErrorCh := make(chan error, 1)

	go func() {
		log.Printf("Server started on host%s\n", server.Addr)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrorCh <- err
		}
	}()

	signalShutdownCh := make(chan os.Signal, 1)
	signal.Notify(signalShutdownCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrorCh:
		log.Fatalf("Server error: %v", err)
	case <-signalShutdownCh:
		log.Println("Resived shutdown signal...")
		t := 5 * time.Second
		ctx, close := context.WithTimeout(context.Background(), t)
		defer close()

		if err := server.Shutdown(ctx); err != nil {
			log.Println("Graceful shutdown failed")
			server.Close()
		}
		log.Println("Graceful shutdown complete")
	}
}
