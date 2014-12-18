// Package message provides ...
package chat

import (
	"time"
)

const (
	_ = iota
	MESSAGE
	DISCONNECT
	JOIN
)

type Message struct {
	Sender  *Client
	Command int
	Message interface{}
	Time    time.Time
}
