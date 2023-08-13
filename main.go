package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	rt := Realtime{Clients: make([]*Client, 0)}

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
	clientID := rt.AddClient(&ch)

	defer close(ch)

	// set the headers
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Starting streaming...\n")
	flusher.Flush()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Client disconnected")
			rt.RemoveClient(clientID)
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
	for _, client := range rt.Clients {
		*client.channel <- msg
	}
	fmt.Fprintf(w, "Message sent: %s", msg)
}

// ///////////////////////////////////////////////////////////////
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

// Represents a single client connection
type Client struct {
	id      uuid.UUID
	channel *chan string
}

// Holds all of the active client connections
type Realtime struct {
	Clients []*Client
}

// Adds a new client to the connection list
func (r *Realtime) AddClient(ch *chan string) uuid.UUID {
	id := uuid.New()
	client := Client{id: id, channel: ch}
	r.Clients = append(r.Clients, &client)
	return id
}

// Removes a client from the connection list
func (r *Realtime) RemoveClient(clientID uuid.UUID) {
	newClients := make([]*Client, 0)

	for _, c := range r.Clients {
		if c.id != clientID {
			newClients = append(newClients, c)
		}
	}

	r.Clients = newClients
}
