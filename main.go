// Package main  provides ...
package main

import (
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	log.Println("New Connection: ", conn.RemoteAddr())
	defer conn.Close()

	buf := make([]byte, 25)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("Remote Closed")
			} else {
				log.Println(conn.RemoteAddr(), "Error: ", err)
			}
			break
		}
		log.Println("Buf", buf, "Recv len:", n)
	}
}

func main() {

	listener, err := net.Listen("tcp", "127.0.0.1:8080")

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleConnection(conn)
	}

}
