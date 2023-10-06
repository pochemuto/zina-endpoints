package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const SECRET_ENV = "ZINA_SECRET"

func handleShutdown(w http.ResponseWriter, r *http.Request) {
	body, error := io.ReadAll(r.Body)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(error.Error()))
		return
	}
	var token = string(body)
	if token != os.Getenv(SECRET_ENV) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf("Bad token: '%s'\n", token)))
		return
	}

	shutdown()
	w.Write([]byte("Shutting down..."))
}

func shutdown() {
	log.Println("Shutting down...")
}

func main() {
	var port = flag.Int("port", 80, "Serving port")
	flag.Parse()

	http.HandleFunc("/shutdown", handleShutdown)

	log.Printf("Listening %d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("Server closed\n")
	} else if err != nil {
		log.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
