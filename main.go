// Package main  provides ...
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const (
	_ = iota
	MESSAGE
	DISCONNECT
	JOIN
)

type Client struct {
	Name      string
	Conn      net.Conn
	Close     bool
	Incomming chan Message
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

func (c *Client) Write(channel chan<- Message) {

	init_msg := &Message{c, JOIN, fmt.Sprintf("%s Joined", c.Name), time.Now()}
	channel <- *init_msg

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
				channel <- *msg
			} else {
				log.Println(c.Conn.RemoteAddr(), "Error: ", err)
				msg = &Message{c, DISCONNECT, fmt.Sprintf("%s DISCONNECT", c.Name), time.Now()}
			}
			break
		} else {
			//log.Println("Buf", buf, "Recv len:", n)
			msg = &Message{c, MESSAGE, line, time.Now()}
		}
		channel <- *msg
	}

}

type Message struct {
	Sender  *Client
	Command int
	Message interface{}
	Time    time.Time
}

func handleConnection(conn net.Conn, global_channel chan Message) *Client {
	log.Println("New Connection: ", conn.RemoteAddr())

	client_incomming := make(chan Message)
	c := &Client{fmt.Sprintf("%s", conn.RemoteAddr()),
		conn, false, client_incomming}

	go c.Listen()
	go c.Write(global_channel)
	return c
}

func distributeMessage(global_channel <-chan Message, clients map[string]*Client) {

	for msg := range global_channel {
		log.Printf("%s@%s:%s", msg.Sender.Name,
			msg.Time.Format(time.RFC3339),
			msg.Message)
		for _, c := range clients {
			c.Incomming <- msg
		}
	}

}

func main() {

	listener, err := net.Listen("tcp", ":12345")

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	global_channel := make(chan Message)
	clients := make(map[string]*Client)
	go distributeMessage(global_channel, clients)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go func(net.Conn, chan Message, map[string]*Client) {
			c := handleConnection(conn, global_channel)
			clients[c.Name] = c
		}(conn, global_channel, clients)
	}

}
