package main

import (
	"fmt"
	"mqtt/shared"
	"net"
)

/*
Each client needs to identify itself with a clientId.
We can then store the list of subscribers as a Dictionary<Topic,List<KVP<ClientID, Conn>>>.
If a connection is closed, we can queue the message for that ClientID for a set amount of time.
*/

// Client - represents a client for our broker.
type Client struct {
	ClientID   string
	Connection net.Conn
}

// Subscription - holds all the clients for a particular topic.
type Subscription struct {
	Topic   string
	Clients []Client
}

func newSubscription(conn net.Conn, clientID, topic string) Subscription {
	return Subscription{Topic: topic, Clients: []Client{Client{ClientID: clientID, Connection: conn}}}
}

func handleMessage(conn net.Conn, subs *[]Subscription) error {
	msg, err := shared.DecodeMessage(conn)
	if err != nil {
		fmt.Println("Failed to decode message: ", err)
		return err
	}

	switch msg.Type {
	case shared.Publish:
		fmt.Println("Publishing message.")
		fmt.Println(subs)
		for _, s := range *subs {
			if s.Topic == msg.Topic {
				fmt.Println("Found subscription.")
				for _, c := range s.Clients {
					err := shared.SendMessage(c.Connection, msg)
					if err != nil {
						fmt.Println("Failed to send message: ", err)
						return err
					}
				}
			}
		}
	case shared.Subscribe:
		fmt.Println("Subscribed client: ", msg.ClientID, conn)
		found := false
		for i := range *subs {
			if (*subs)[i].Topic == msg.Topic {
				(*subs)[i].Clients = append((*subs)[i].Clients, Client{ClientID: msg.ClientID, Connection: conn})
				found = true
			}
		}
		if !found {
			*subs = append(*subs, newSubscription(conn, msg.ClientID, msg.Topic))
		}
		fmt.Println(subs)
	}
	return nil
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	subs := []Subscription{}

	fmt.Println("Successfully created tcp listener.")
	for {
		fmt.Println("\nWaiting for connection...")
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed accept connection: ", err)
			return
		}
		go handleMessage(conn, &subs)
	}
}
