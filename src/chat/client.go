// Package client provides ...
package chat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type Client struct {
	Server *ChatServer
	Name   string
	Conn   net.Conn
	Rooms  map[string]*Room
	In     chan *Message
	Out    chan *Message
	Quit   chan bool
}

func (c *Client) Listen() {
	log.Printf("New client: %s", c.Name)
	for {
		select {
		case msg := <-c.In:
			// Client receive message
			go c.Write(msg)
		case msg := <-c.Out:
			switch msg.Command {
			case QUIT:
				// broadcast to all rooms
				for _, r := range c.Rooms {
					r.In <- msg
				}
				c.Quit <- true
			case JOIN:
				name := msg.Receiver
				var room *Room
				if _, ok := c.Server.Rooms[name]; !ok {
					// Not exist room make one.
					room = &Room{c.Server,
						name,
						make(map[string]*Client, 0),
						make(chan *Message),
						make(chan bool)}
					c.Server.Rooms[name] = room
					go room.Listen()
				}

				c.Rooms[name] = c.Server.Rooms[name]
				c.Rooms[name].In <- msg
			default:
				c.Rooms[msg.Receiver].In <- msg
			}
		case <-c.Quit:
			return
		}
	}
}

func (c *Client) Write(msg *Message) {

	c.Conn.Write([]byte(fmt.Sprintf("%s %s:%s\n",
		msg.Time.Format(time.RFC3339),
		msg.Sender.Name,
		msg.Content)))
}

func (c *Client) Recv() {

	buf := bufio.NewReader(c.Conn)
	var msg *Message

	for {
		line, err := buf.ReadString('\n')

		if err != nil || len(line) == 0 {
			if err == io.EOF || len(line) == 0 {
				log.Println(c.Name, " Remote Closed")
				msg = &Message{c, "", QUIT, fmt.Sprintf("%s Lefted", c.Name), time.Now()}
			} else {
				log.Println(c.Conn.RemoteAddr(), "Error: ", err)
				msg = &Message{c, "", QUIT, fmt.Sprintf("%s DISCONNECT", c.Name), time.Now()}
			}
			c.Out <- msg
			break
		} else {

			data := strings.Split(strings.TrimSpace(line), " ")

			room, content := data[0], data[1]
			if _, ok := c.Rooms[room]; ok {
				msg = &Message{c, room, NORMAL, content, time.Now()}
			} else {
				msg = &Message{c, room, JOIN, content, time.Now()}
			}
		}
		c.Out <- msg
	}

}
