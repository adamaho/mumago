package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	ch := make(chan string, 10)
	defer close(ch)

	r.Group(func(r chi.Router) {
		r.Get("/subscribe", func(w http.ResponseWriter, r *http.Request) { subscribe(w, r, ch) })
		r.Post("/publish/{message}", func(w http.ResponseWriter, r *http.Request) { publish(w, r, ch) })
	})

	log.Println("starting server at https://localhost:3000")

	err := http.ListenAndServeTLS("localhost:3000", "cert.pem", "key.pem", r)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func subscribe(w http.ResponseWriter, r *http.Request, c chan string) {
	// get the channel from context
	ctx := r.Context()

	// check if the user supports http2
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming not supported!", http.StatusInternalServerError)
		return
	}

	// set the headers
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Starting streaming...\n")
	flusher.Flush()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Client disconnected")
			return
		case value := <-c:
			if value != "" {
				fmt.Fprintf(w, "data: %s\n", value)
			}
			flusher.Flush()
		}
	}
}

func publish(w http.ResponseWriter, r *http.Request, c chan string) {
	msg := chi.URLParam(r, "message")
	c <- msg
	fmt.Fprintf(w, "Message sent: %s", msg)
}

// /////////////////////////////////////////////////////////////////
func sleep(t time.Duration) {
	time.Sleep(t * time.Second)
}

type Payload struct {
	Data []byte
}

// Creates a new plain text payload
func NewPayload(s string) *Payload {
	return &Payload{Data: []byte(s)}
}

// Creates a new json payload
func NewJsonPayload(j struct{}) (*Payload, error) {
	d, err := json.Marshal(j)

	if err != nil {
		log.Println("failed to marshal struct:", err)
		return nil, err
	}

	return &Payload{Data: d}, nil
}

// Represents the client that recieves messages in the channel
type Sender struct {
	rx chan<- Payload
}

// Holds all of the active client connections
type JsonPatchStream struct {
	Clients []Sender
}

func (j JsonPatchStream) Subscribe() string {
	return "Streammmmmm"
}
