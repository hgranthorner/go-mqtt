package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Successfully created tcp listener.")
	for {
		fmt.Println("Waiting for connection...")
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		for {
			bytes := make([]byte, 2048)
			n, err := conn.Read(bytes)
			if err != nil && err != io.EOF {
				fmt.Println("Error: ", err)
				break
			}
			if n > 0 {
				fmt.Println("n: ", n)
				fmt.Println(string(bytes))
			}
		}
	}
}
