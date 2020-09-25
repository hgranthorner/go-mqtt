package shared

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

// MessageType - go-enum of different options in the Type field of the Message type.
type MessageType int

const (
	// Publish - indicates that the client is sending data on the stream.
	// The Payload should not be nil in this case.
	Publish MessageType = iota
	// Subscribe - indicates that the client would like to be notified of messages on a certain topic.
	// The Payload field should be nil in this case.
	Subscribe
	EndMessage
)

func (mt MessageType) String() string {
	return [...]string{"Publish", "Subscribe", "EndMessage"}[mt]
}

// Message - contains the sender's key, the deliverer's key, and a JSON serialized message.
type Message struct {
	Type     MessageType
	ClientID string
	Topic    string
	Payload  string
}

// DecodeMessage - attempts to decode a Message from a TCP connection.
func DecodeMessage(conn net.Conn) (Message, error) {
	var msg Message
	bytes := make([]byte, 2048)
	n, err := conn.Read(bytes)
	if err != nil && err != io.EOF {
		return msg, err
	}

	// Should be replaced with gob package in standard library
	err = json.Unmarshal(bytes[0:n], &msg)

	if err != nil {
		fmt.Println(string(bytes))

		return msg, err
	}
	return msg, nil
}

// SendMessage - send a Message over a given Conn.
func SendMessage(conn net.Conn, msg Message) error {
	// Should be replaced with gob package in standard library
	str, err := json.Marshal(msg)

	if err != nil {
		fmt.Println("Failed to marshal: ", err)
		return err
	}
	_, err = conn.Write(str)
	if err != nil {
		fmt.Println("Failed to write to TCP: ", err)
		return err
	}

	return nil
}
