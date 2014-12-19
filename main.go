// Package main  provides ...
package main

import (
	"chat"
)

func main() {
	server := &chat.ChatServer{":12345", make(map[string]*chat.Room)}
	server.ListenAndServe()
}
