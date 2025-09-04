package thebus

import "time"

type Message struct {
	Topic     string
	Timestamp time.Time
	Payload   []byte
	Seq       uint64
}
