// Package main provides ...
package chat

import (
	"fmt"
	"log"
)

func kickName(msg *Message) string {
	return fmt.Sprintf("%s", msg.Message)
}

type Room struct {
	Name    string
	clients map[string]*Client
	In      chan *Message
	Quit    chan bool
}

func (r *Room) Listen() {

	log.Printf("Chatroom: ", r.Name, " opened")
	for {
		select {
		case msg := <-r.In:
			switch msg.Command {
			case QUIT:
				delete(r.clients, msg.Sender.Name)
				go r.broadcast(msg)
			case JOIN:
				r.clients[msg.Sender.Name] = msg.Sender
				go r.broadcast(msg)
			case KICK:
				name := kickName(msg)
				if val, ok := r.clients[name]; ok {
					delete(r.clients, name)
					go r.broadcast(msg)
				}
			case DISMISS:
				// Blocking broadcasting...
				r.broadcast(msg)
				r.Quit <- true
			default:
				go r.broadcast(msg)
			}

		case <-r.Quit:
			for k := range r.clients {
				delete(r.clients, k)
			}
			log.Printf("Chatroom: ", r.Name, " closed")
			// Ok, Mission Completed
			return
		}
	}
}

func (r *Room) broadcast(msg *Message) {

	for _, c := range r.clients {
		c.Incomming <- msg
	}

}
