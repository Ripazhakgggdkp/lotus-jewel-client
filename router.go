package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func listenHTTP(client *Client) <-chan error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	route(r, client)

	server := &http.Server{Addr: ":25565", Handler: r}

	error := make(chan error)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			error <- err
		}
	}()

	// Handle shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func(chan os.Signal) {
		<-stop
		log.Println("Shutting down client")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Could not shutdown server gracefully", err)
		}
		os.Exit(0)
	}(stop)

	return error

}

// Handle requests
func route(r *chi.Mux, client *Client) {
	r.Get("/vibrate/{strength}", func(w http.ResponseWriter, r *http.Request) {
		strength, err := strconv.ParseFloat(chi.URLParam(r, "strength"), 64)

		if err != nil {
			strength = 0
		}

		client.vibrate(strength)
	})
	r.Get("/connect", func(w http.ResponseWriter, r *http.Request) {
		err := client.connect(1)
		if err != nil {
			log.Fatal("Could not connect game to Intiface Central", err)
		}
	})
	r.Get("/stop", func(w http.ResponseWriter, r *http.Request) {
		client.stop()
	})
}
