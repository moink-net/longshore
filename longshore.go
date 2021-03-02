package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

func getDockerContainerStats() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		// set status docker_client_setup_healthy: false
		return
	}
	client.SkipServerVersionCheck = true
}

var (
	httpListenAddr string
)

func main() {
	log.SetFlags(0)
	log.SetOutput(Iso8601Writer{os.Stdout})
	log.Println("Starting `longshore`...")

	flag.StringVar(&httpListenAddr, "http-listen-addr", ":8080", "http server listen address")
	flag.Parse()

	httpRouter := http.NewServeMux()
	httpRouter.Handle("/", daemonStats())
	httpRouter.Handle("/livez", livez())
	httpRouter.Handle("/readyz", readyz())
	httpRouter.Handle("/healthz", healthz())

	server := &http.Server{
		Addr:         httpListenAddr,
		Handler:      httpRouter,
		ErrorLog:     log.New(Iso8601Writer{os.Stdout}, "http: ", 0),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// handle graceful shutdown
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("Stopping `longshore`...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("could not gracefully stop: %v\n", err)
		}
		close(done)
	}()

	// start http server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on %s: %v\n", httpListenAddr, err)
	}

}

type jsonError struct {
	Err string `json:"error"`
}

func daemonStats() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// create Docker client
		client, err := docker.NewClientFromEnv()
		client.SkipServerVersionCheck = true

		// unhealthy; unable to create Docker client
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			jsonOut, _ := json.Marshal(jsonError{err.Error()})
			w.Write(jsonOut)
			return
		}

		// get daemon info
		daemonInfo, err := client.Info()

		// unhealthy; unable to get Docker daemon info
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			jsonOut, _ := json.Marshal(jsonError{err.Error()})
			w.Write(jsonOut)
			return
		}

		w.WriteHeader(http.StatusOK)
		jsonOut, _ := json.Marshal(daemonInfo)
		w.Write(jsonOut)
		return
	})
}

func livez() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		return
	})
}
func readyz() http.Handler {
	return daemonStats()
}
func healthz() http.Handler {
	return readyz()
}
