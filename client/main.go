package main

import (
	"fmt"
	"mqtt/shared"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Failed to dial to port: ", err)
	}

	if os.Args[1] == "p" {
		msg := shared.Message{Type: shared.Publish, ClientID: "me", Topic: "you", Payload: "{ \"id\": 1 }"}
		shared.SendMessage(conn, msg)
	} else if os.Args[1] == "s" {
		msg := shared.Message{Type: shared.Subscribe, ClientID: "me", Topic: "you", Payload: ""}
		shared.SendMessage(conn, msg)
		// shared.SendMessage(conn, shared.Message{Type: shared.EndMessage, ClientID: "", Topic: "", Payload: ""})
		fmt.Println("Subscribing to topic and waiting for message...")
		for {
			msg, err := shared.DecodeMessage(conn)
			if err != nil {
				fmt.Println("Failed to decode message: ", err)
				return
			}
			fmt.Println(msg)
		}
	} else if os.Args[1] == "t" {

	} else {
		fmt.Println("Program requires one argument: p for publish, s for subscribe")
	}
}
