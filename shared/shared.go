package shared

// Message - contains the sender's key, the deliverer's key, and a JSON serialized message.
type Message struct {
	SenderKey   string
	DeliveryKey string
	Payload     string
}

// IDPayload - basic payload type with just an integer id.
type IDPayload struct {
	ID int
}
