package thebus

import "time"

type Message struct {
	Topic     string
	Timestamp time.Time
	Payload   []byte
	Seq       uint64
}

type messageRef struct {
	topic   string
	ts      time.Time
	seq     uint64
	payload []byte
}

func makeMessage(topic string, mr messageRef, sub *subscription) Message {
	msg := Message{
		Topic:     topic,
		Timestamp: mr.ts,
		Seq:       mr.seq,
	}
	if sub.cfg.Strategy == SubscriptionStrategyPayloadShared {
		msg.Payload = mr.payload
	} else {
		buf := make([]byte, len(mr.payload))
		copy(buf, mr.payload)
		msg.Payload = buf
	}
	return msg
}
