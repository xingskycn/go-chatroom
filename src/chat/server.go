// Package server provides ...
package chat

import (
	"log"
	"net"
)

func distributeMessage(global_channel <-chan *Message, clients map[string]*Client) {

	for msg := range global_channel {
		//log.Printf("%s@%s: %s", msg.Sender.Name,
		//	msg.Time.Format(time.RFC3339),
		//	msg.Message)
		if msg.Command == DISCONNECT {
			msg.Sender.Quit <- true
			delete(clients, msg.Sender.Name)
		}
		go func(clients map[string]*Client) {
			for _, c := range clients {
				c.Incomming <- msg
			}
		}(clients)
	}

}

func Serve() {

	listener, err := net.Listen("tcp", ":12345")

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	global_channel := make(chan *Message)
	clients := make(map[string]*Client)

	go distributeMessage(global_channel, clients)

	// Main loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go func(net.Conn, chan *Message, map[string]*Client) {
			c := handleConnection(conn, global_channel)
			clients[c.Name] = c
		}(conn, global_channel, clients)
	}

}
