package shared

import (
	"encoding/binary"
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
)

func (mt MessageType) String() string {
	return [...]string{"Publish", "Subscribe"}[mt]
}

// Message - contains the sender's key, the deliverer's key, and a JSON serialized message.
type Message struct {
	Type     MessageType
	ClientID string
	Topic    string
	Payload  string
}

// DecodeMessage - attempts to decode a Message from a TCP connection.
func DecodeMessage(conn net.Conn) ([]Message, error) {
	var msg Message
	msgs := []Message{}

	bytes := make([]byte, 2048)
	count, err := conn.Read(bytes)
	if err != nil && err != io.EOF {
		return msgs, err
	}

	start := 0
	total := 0
	for {
		start = total
		size, _ := binary.Varint(bytes[start : start+32])
		start += 32

		// Should be replaced with gob package in standard library
		err = json.Unmarshal(bytes[start:(start+int(size))], &msg)
		total = start + int(size)

		if err != nil {
			fmt.Println(string(bytes[32:(32 + size)]))

			return msgs, err
		}
		msgs = append(msgs, msg)
		if total >= count {
			break
		}
	}
	return msgs, nil
}

// SendMessage - send a Message over a given Conn.
func SendMessage(conn net.Conn, msg Message) error {
	// Should be replaced with gob package in standard library
	str, err := json.Marshal(msg)

	if err != nil {
		fmt.Println("Failed to marshal: ", err)
		return err
	}
	bytes := make([]byte, 32)

	binary.PutVarint(bytes, int64(len(str)))
	size, _ := binary.Varint(bytes)
	fmt.Println(size)

	_, err = conn.Write(bytes)

	_, err = conn.Write(str)
	if err != nil {
		fmt.Println("Failed to write to TCP: ", err)
		return err
	}

	return nil
}
