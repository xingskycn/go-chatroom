// Package main provides ...
package chat

import (
	"fmt"
	"log"
)

type Room struct {
	clients   map[string]*Client
	Incomming chan *Message
}

func (r *Room) Listen() {

	for {
		select {
		case msg := <-r.Incomming:
			log.Printf(fmt.Sprintf("%s", msg.Message))
			switch msg.Command {
			case QUIT:
				delete(r.clients, msg.Sender.Name)
			case JOIN:
				r.clients[msg.Sender.Name] = msg.Sender
			}
			go r.Broadcast(msg)
		}
	}
}

func (r *Room) Broadcast(msg *Message) {

	for _, c := range r.clients {
		c.Incomming <- msg
	}

}
