// Package message provides ...
package chat

import (
	"time"
)

const (
	_ = iota
	NORMAL
	QUIT
	JOIN
	DISMISS
	PAUSE
	KICK
)

type Message struct {
	Sender  *Client
	Command int
	Message interface{}
	Time    time.Time
}
