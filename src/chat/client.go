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
	Name      string
	Conn      net.Conn
	Close     bool
	Incomming chan *Message
	Quit      chan bool
}

func (c *Client) Listen() {

	for msg := range c.Incomming {
		if c.Close {
			break
		}
		c.Conn.Write([]byte(fmt.Sprintf("%s %s:%s\n",
			msg.Time.Format(time.RFC3339),
			msg.Sender.Name,
			msg.Message)))
	}
}

func (c *Client) Write(channel chan<- *Message) {

	init_msg := &Message{c, JOIN, fmt.Sprintf("%s Joined", c.Name), time.Now()}
	channel <- init_msg

	buf := bufio.NewReader(c.Conn)

	for {
		if c.Close {
			break
		}
		var msg *Message
		line, err := buf.ReadString('\n')
		if err != nil || len(line) == 0 {
			if err == io.EOF || len(line) == 0 {
				log.Println(c.Name, " Remote Closed")
				msg = &Message{c, DISCONNECT, fmt.Sprintf("%s Lefted", c.Name), time.Now()}
			} else {
				log.Println(c.Conn.RemoteAddr(), "Error: ", err)
				msg = &Message{c, DISCONNECT, fmt.Sprintf("%s DISCONNECT", c.Name), time.Now()}
			}
			channel <- msg
			break
		} else {
			//log.Println("Buf", buf, "Recv len:", n)
			msg = &Message{c, MESSAGE, line, time.Now()}
		}
		channel <- msg
	}

}

func (c *Client) ListenQuit() {
	c.Close = <-c.Quit
}

func handleConnection(conn net.Conn, global_channel chan<- *Message) *Client {
	log.Println("New Connection: ", conn.RemoteAddr())

	client_incomming := make(chan *Message)

	buf := bufio.NewReader(conn)
	text := []byte("What is your name?\n")
	conn.Write(text)

	line, err := buf.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	c := &Client{fmt.Sprintf("%s", strings.TrimSpace(line)),
		conn, false, client_incomming, make(chan bool)}

	go c.Listen()
	go c.Write(global_channel)
	go c.ListenQuit()
	return c
}
