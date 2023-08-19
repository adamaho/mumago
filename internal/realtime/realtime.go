package realtime

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mattbaird/jsonpatch"
)

// Represents a single client connection
type Client struct {
	id      uuid.UUID
	Channel *chan []byte
}

// Holds all of the active client connections
type Realtime struct {
	Clients []*Client
	Data    json.RawMessage
}

// Creates a new instance of realtime
func New() Realtime {
	return Realtime{Clients: make([]*Client, 0)}
}

// Handles creating a stream and channel if the X-Muma-Stream header is set. If header
// is not set, the raw json message is returned to the user like a standard REST api.
func (rt *Realtime) Stream(w http.ResponseWriter, r *http.Request, d json.RawMessage) {
	ctx := r.Context()

	s := r.Header.Get("X-Muma-Stream")
	if s == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(d)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		// TODO: come up with a generic error struct to return instead of a plain string
		http.Error(w, "Streaming not supported!", http.StatusInternalServerError)
		return
	}

	ch := make(chan []byte, 10)
	defer close(ch)

	clientID := rt.AddClient(&ch)

	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Content-Type", "application/json+ndjsonpatch")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", d)
	flusher.Flush()

	rt.Data = d

	for {
		select {
		case <-ctx.Done():
			rt.RemoveClient(clientID)
			return
		case value := <-ch:
			if len(value) > 0 {
				fmt.Fprintf(w, "%s\n", value)
			}
			flusher.Flush()
		}
	}
}

// Adds a new client to the connection list
func (rt *Realtime) AddClient(ch *chan []byte) uuid.UUID {
	id := uuid.New()
	client := Client{id: id, Channel: ch}
	rt.Clients = append(rt.Clients, &client)
	return id
}

// Removes a client from the connection list
func (rt *Realtime) RemoveClient(clientID uuid.UUID) {
	newClients := make([]*Client, 0)

	for _, c := range rt.Clients {
		if c.id != clientID {
			newClients = append(newClients, c)
		}
	}

	rt.Clients = newClients
}

// Creates a new json patch
func (rt *Realtime) PublishPatch(target json.RawMessage) {
	patch, _ := jsonpatch.CreatePatch(rt.Data, target)
	patchJson, err := json.Marshal(patch)

	if err != nil {
		log.Print("Failed to marshal json for patch")
		return
	}

	for _, client := range rt.Clients {
		*client.Channel <- patchJson
	}

	rt.Data = target
}
