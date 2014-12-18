// Package main provides ...
package chat

type Room struct {
	clients   map[string]*Client
	Incomming chan *Message
}

func (r *Room) Listen() {
	for msg := range r.Incomming {

	}
}
