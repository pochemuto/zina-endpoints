package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
)

const SECRET_ENV = "ZINA_SECRET"

var app_version string

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func doResponse(w http.ResponseWriter, message string) {
	w.Write([]byte(message))
	log.Println(message)
}

func handleShutdown(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	var token = string(body)
	if token != os.Getenv(SECRET_ENV) {
		w.WriteHeader(http.StatusForbidden)
		doResponse(w, fmt.Sprintf("Bad token: '%s'\n", token))
		return
	}

	err = shutdown()
	if err != nil {
		handleError(w, err)
		return
	}
	doResponse(w, "Shutting down...")
}

func shutdown() error {
	wr, err := os.OpenFile("/endpoints-pipe", os.O_WRONLY|syscall.O_NONBLOCK, 0o600)
	if err != nil {
		return err
	}
	defer wr.Close()
	if _, err := wr.Write([]byte("shutdown")); err != nil {
		return err
	}
	log.Println("Shutting down...")
	return nil
}

func main() {
	var port = flag.Int("port", 80, "Serving port")
	flag.Parse()

	http.HandleFunc("/shutdown", handleShutdown)

	log.Printf("Version %s\n", app_version)
	log.Printf("Listening :%d ...\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("Server closed\n")
	} else if err != nil {
		log.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
