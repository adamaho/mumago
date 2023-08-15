package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"mumago/internal/realtime"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	rt := realtime.Realtime{Clients: make([]*realtime.Client, 0)}

	r.Group(func(r chi.Router) {
		r.Get("/subscribe", func(w http.ResponseWriter, r *http.Request) { subscribe(w, r, &rt) })
		r.Post("/publish/{message}", func(w http.ResponseWriter, r *http.Request) { publish(w, r, &rt) })
	})

	log.Println("starting todos server at https://localhost:3000")

	err := http.ListenAndServeTLS("localhost:3000", "cert.pem", "key.pem", r)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func subscribe(w http.ResponseWriter, r *http.Request, rt *realtime.Realtime) {
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
			fmt.Printf("Client disconnected\n")
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

func publish(w http.ResponseWriter, r *http.Request, rt *realtime.Realtime) {
	msg := chi.URLParam(r, "message")
	for _, client := range rt.Clients {
		*client.Channel <- msg
	}
	fmt.Fprintf(w, "Message sent: %s", msg)
}
