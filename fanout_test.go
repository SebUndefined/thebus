package thebus

import (
	"testing"
	"time"
)

func TestSnapshotSubsLocked(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]*subscription
		output []*subscription
	}{
		{
			name:   "empty",
			input:  make(map[string]*subscription),
			output: make([]*subscription, 0),
		},
		{
			name: "found",
			input: map[string]*subscription{
				"ID1": {
					subscriptionID:  "ID1",
					cfg:             SubscriptionConfig{},
					topic:           "",
					messageChan:     nil,
					messages:        nil,
					unsubscribeFunc: nil,
				},
			},
			output: []*subscription{
				{
					subscriptionID:  "ID1",
					cfg:             SubscriptionConfig{},
					topic:           "",
					messageChan:     nil,
					messages:        nil,
					unsubscribeFunc: nil,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := snapshotSubsLocked(test.input)
			if len(test.output) != len(output) {
				t.Errorf("len(input) = %d, want %d", len(test.input), len(test.output))
			}
		})
	}
}

type testInputTryDelivery struct {
	sub   *subscription
	msg   Message
	timer *time.Timer
}

func TestTryDeliver(t *testing.T) {
	tests := []struct {
		name   string
		input  testInputTryDelivery
		output bool
	}{
		{
			name: "success",
			input: testInputTryDelivery{
				&subscription{},
				Message{},
				nil,
			},
			output: true,
		},
		{
			name: "failure",
			input: testInputTryDelivery{
				&subscription{},
				Message{},
				nil,
			},
			output: false,
		},
		{
			name: "success_timeout",
			input: testInputTryDelivery{
				&subscription{},
				Message{},
				nil,
			},
			output: true,
		},
		{
			name: "failure_timeout",
			input: testInputTryDelivery{
				&subscription{},
				Message{},
				nil,
			},
			output: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
