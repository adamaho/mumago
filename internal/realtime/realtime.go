package realtime

import "github.com/google/uuid"

// Represents a single client connection
type Client struct {
	id      uuid.UUID
	Channel *chan string
}

// Holds all of the active client connections
type Realtime struct {
	Clients []*Client
}

// TODO: Create New function for creating a realtime struct
// TODO: Create Streaming function for the streaming endpoint
// TODO: It should handle non-streaming support as well and streaming should be supported based on a header
// TODO: Should have a mutation function maybe?

// Adds a new client to the connection list
func (r *Realtime) AddClient(ch *chan string) uuid.UUID {
	id := uuid.New()
	client := Client{id: id, Channel: ch}
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
