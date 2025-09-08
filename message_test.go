package thebus

import (
	"testing"
	"time"
)

func TestMakeMessage(t *testing.T) {
	tests := []struct {
		name       string
		topic      string
		messageRef messageRef
		sub        *subscription
		outPut     Message
	}{
		{
			name:  "Payload Shared",
			topic: "shared",
			sub: &subscription{
				cfg: SubscriptionConfig{
					Strategy: SubscriptionStrategyPayloadShared,
				},
			},
			messageRef: messageRef{
				topic:   "shared",
				ts:      time.Now(),
				seq:     0,
				payload: []byte("payload"),
			},
		},
		{
			name:  "Payload cloned",
			topic: "cloned",
			sub: &subscription{
				cfg: SubscriptionConfig{
					Strategy: SubscriptionStrategyPayloadShared,
				},
			},
			messageRef: messageRef{
				topic:   "cloned",
				ts:      time.Now(),
				seq:     0,
				payload: []byte("payload"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := makeMessage(test.topic, test.messageRef, test.sub)
			var sameBacking bool
			if len(m.Payload) > 0 && len(test.messageRef.payload) > 0 {
				if &m.Payload[0] == &test.messageRef.payload[0] {
					sameBacking = true
				}
			}
			if test.sub.cfg.Strategy == SubscriptionStrategyPayloadShared && !sameBacking {
				t.Fatal("Expected payload to be shared")
			} else if test.sub.cfg.Strategy == SubscriptionStrategyPayloadClonedPerSubscriber && sameBacking {
				t.Fatal("Expected payload to be cloned")
			}
		})
	}
}
