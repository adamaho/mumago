package realtime

type Payload struct {
	Data []byte
}

// Creates a new plain text payload
func NewPayload(s string) *Payload {
	return &Payload{Data: []byte(s)}
}
