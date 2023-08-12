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

	rt := Realtime{Clients: make([]*chan string, 0)}

	r.Group(func(r chi.Router) {
		r.Get("/subscribe", func(w http.ResponseWriter, r *http.Request) { subscribe(w, r, &rt) })
		r.Post("/publish/{message}", func(w http.ResponseWriter, r *http.Request) { publish(w, r, &rt) })
	})

	log.Println("starting server at https://localhost:3000")

	err := http.ListenAndServeTLS("localhost:3000", "cert.pem", "key.pem", r)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func subscribe(w http.ResponseWriter, r *http.Request, rt *Realtime) {
	// init the context
	ctx := r.Context()

	// check if the user supports http2
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming not supported!", http.StatusInternalServerError)
		return
	}

	// create channel and add to list of clients
	ch := make(chan string, 10)
	rt.AddClient(&ch)

	defer close(ch)

	// set the headers
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Starting streaming...\n")
	flusher.Flush()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Client disconnected")
			// TODO: Remove client from list
			// TODO: Create unique id for each client to make it easier to remove
			// TODO: Create Client struct which holds an id and a pointer to a channel
			return
		case value := <-ch:
			if value != "" {
				fmt.Fprintf(w, "data: %s\n", value)
			}
			flusher.Flush()
		}
	}
}

func publish(w http.ResponseWriter, r *http.Request, rt *Realtime) {
	msg := chi.URLParam(r, "message")
	for _, ch := range rt.Clients {
		*ch <- msg
	}
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

// Holds all of the active client connections
type Realtime struct {
	Clients []*chan string
}

// Adds a new client to the realtime sync
func (r *Realtime) AddClient(ch *chan string) {
	r.Clients = append(r.Clients, ch)
}
