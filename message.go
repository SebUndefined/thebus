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

func newMessageRef(topic string, ts time.Time, seq uint64, payload []byte) messageRef {
	return messageRef{
		topic:   topic,
		ts:      ts,
		seq:     seq,
		payload: payload,
	}
}

func newMessageRefCopy(topic string, ts time.Time, seq uint64, payload []byte) messageRef {
	newPayload := make([]byte, len(payload))
	copy(newPayload, payload)
	return messageRef{
		topic:   topic,
		ts:      ts,
		seq:     seq,
		payload: newPayload,
	}
}
