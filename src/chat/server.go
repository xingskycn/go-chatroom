// Package server provides ...
package chat

import (
	"fmt"
	"log"
	"net"
	"time"
)

type ChatServer struct {
	Bind_to string
	Rooms   map[string]*Room
}

func (server *ChatServer) reportStatus() {

	for {
		time.Sleep(10 * time.Second)
		for _, room := range server.Rooms {
			fmt.Printf("%s:%d", room.Name, len(room.Clients))
		}
	}

}

func (server *ChatServer) ListenAndServe() {

	listener, err := net.Listen("tcp", server.Bind_to)

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	go server.reportStatus()
	// Main loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go func(conn net.Conn, server *ChatServer) {
			c := &Client{server, fmt.Sprintf("%s", conn.RemoteAddr()),
				conn, make(map[string]*Room), make(chan *Message),
				make(chan *Message), make(chan bool)}
			go c.Listen()
			go c.Recv()
		}(conn, server)
	}

}
