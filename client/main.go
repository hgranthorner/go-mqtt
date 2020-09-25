package main

import (
	"encoding/json"
	"fmt"
	"mqtt/shared"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}

	msg := shared.Message{SenderKey: "me", DeliveryKey: "you", Payload: "{\"id\": 1}"}

	str, err := json.Marshal(msg)

	if err != nil {
		fmt.Println("Failed to marshal: ", err)
		return
	}
	conn.Write(str)
	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	//conn.Write([]byte("Test"))
	//time.Sleep(1 * time.Second)
	//conn.Write([]byte("Test 2"))
}
