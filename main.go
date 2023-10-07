package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const SECRET_ENV = "ZINA_SECRET"

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
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
		w.Write([]byte(fmt.Sprintf("Bad token: '%s'\n", token)))
		return
	}

	err = shutdown()
	if err != nil {
		handleError(w, err)
		return
	}
	w.Write([]byte("Shutting down..."))
}

func shutdown() error {
	cmd := exec.Command("systemctl", "poweroff")
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Println("Shutting down...")
	return nil
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
