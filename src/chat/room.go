// Package main provides ...
package chat

import (
	"fmt"
	"log"
)

func kickName(msg *Message) string {
	return fmt.Sprintf("%s", msg.Content)
}

type Room struct {
	Server  *ChatServer
	Name    string
	Clients map[string]*Client
	In      chan *Message
	Quit    chan bool
}

func (r *Room) Listen() {

	log.Printf("Chatroom: %s opened", r.Name)
	for {
		select {
		case msg := <-r.In:
			switch msg.Command {
			case QUIT:
				delete(r.Clients, msg.Sender.Name)
				go r.broadcast(msg)
			case JOIN:
				log.Printf("%s joined", msg.Sender.Name)
				r.Clients[msg.Sender.Name] = msg.Sender
				go r.broadcast(msg)
			case KICK:
				name := kickName(msg)
				if _, ok := r.Clients[name]; ok {
					delete(r.Clients, name)
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
			delete(r.Server.Rooms, r.Name)
			for k := range r.Clients {
				delete(r.Clients, k)
			}
			log.Printf("Chatroom: %s closed", r.Name)
			// Ok, Mission Completed
			return
		}
	}
}

func (r *Room) broadcast(msg *Message) {

	for _, c := range r.Clients {
		c.In <- msg
	}

}
