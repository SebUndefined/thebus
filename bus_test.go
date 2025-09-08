package thebus

import (
	"context"
	"testing"
	"time"
)

func TestWithReadState(t *testing.T) {
	b, _ := New()
	bb := b.(*bus)

	err := bb.withReadState("nope", func(st *topicState) error {
		if st != nil {
			t.Fatal("expected nil state for unknown topic")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// Test if exists
	err = bb.withWriteState("nope", true, func(st *topicState) error {
		if st == nil {
			t.Fatal("expected non-nil state for unknown topic")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithWriteStateNoCreation(t *testing.T) {
	b, _ := New()
	bb := b.(*bus)

	called := false
	err := bb.withWriteState("t", false, func(st *topicState) error {
		called = true
		if st != nil {
			t.Fatal("expected nil state when createIfNotExists=false and topic missing")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("callback not called")
	}
}

func TestWithWriteStateCreateIfNotExists(t *testing.T) {
	b, _ := New()
	bb := b.(*bus)

	err := bb.withWriteState("t", true, func(st *topicState) error {
		if st == nil {
			t.Fatal("state should be created")
		}
		if st.inQueue == nil {
			t.Fatal("inQueue should be initialized")
		}
		if !st.started.Load() {
			t.Fatal("worker should be marked started")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
func TestWithWriteStateCreateMaxTopics(t *testing.T) {
	b, _ := New(WithMaxTopics(1))
	defer b.Close()
	bb := b.(*bus)

	// 1rst topic should be ok
	if err := bb.withWriteState("t1", true, func(st *topicState) error {
		if st == nil {
			t.Fatal("st should not be nil")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	// use same topic, should not block
	if err := bb.withWriteState("t1", true, func(st *topicState) error {
		if st == nil {
			t.Fatal("st should not be nil (existing topic)")
		}
		return nil
	}); err != nil {
		t.Fatalf("unexpected err on existing topic: %v", err)
	}

	// Break ! max topic should border the creation
	if err := bb.withWriteState("t2", true, func(st *topicState) error { return nil }); err == nil {
		t.Fatal("expected error MaxTopics exceeded")
	}

	// The topic 1 is stated (just for clarification
	bb.mutex.RLock()
	st := bb.subscriptions["t1"]
	bb.mutex.RUnlock()
	if st == nil || !st.started.Load() {
		t.Fatal("topic t1 should exist and be started")
	}
}

func TestWithWriteStateCreateStartTopicRead(t *testing.T) {
	b, _ := New()
	bb := b.(*bus)

	err := bb.withWriteState("t", true, func(st *topicState) error {
		// Add a message for checking the health
		select {
		case st.inQueue <- messageRef{topic: "t", ts: time.Now(), payload: []byte("x")}:
		default:
			t.Fatal("inQueue should be writable")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Wait for worker to read, should not panic
	time.Sleep(10 * time.Millisecond)
}

func TestSnapshotSubsLocked(t *testing.T) {
	m := map[string]*subscription{
		"A": {subscriptionID: "A"},
	}
	out := snapshotSubsLocked(m)
	if len(out) != 1 {
		t.Fatalf("want 1, got %d", len(out))
	}

	// Mutate the map, should not modify the main map from the bus
	delete(m, "A")
	if len(out) != 1 {
		t.Fatal("snapshot should be independent from map mutation")
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
				&subscription{
					subscriptionID: "ID1",
					cfg: SubscriptionConfig{
						DropIfFull: true,
					},
					topic:       "test",
					messageChan: make(chan Message, 1),
				},
				Message{},
				nil,
			},
			output: true,
		},
		{
			name: "failure",
			input: testInputTryDelivery{
				&subscription{
					subscriptionID: "ID1",
					cfg: SubscriptionConfig{
						DropIfFull: true,
					},
				},
				Message{
					Topic:     "test",
					Timestamp: time.Now(),
					Payload:   []byte("test"),
					Seq:       0,
				},
				nil,
			},
			output: false,
		},
		{
			name: "success_timeout",
			input: testInputTryDelivery{
				&subscription{
					subscriptionID: "ID1",
					cfg: SubscriptionConfig{
						DropIfFull:  false,
						SendTimeout: 2 * time.Second,
					},
					topic:       "test",
					messageChan: make(chan Message, 1),
				},
				Message{},
				time.NewTimer(2 * time.Second),
			},
			output: true,
		},
		{
			name: "failure_timeout",
			input: testInputTryDelivery{
				&subscription{
					subscriptionID: "ID2",
					cfg: SubscriptionConfig{
						DropIfFull:  false,
						SendTimeout: 50 * time.Millisecond, // court timeout
					},
					topic:       "test",
					messageChan: make(chan Message, 1),
				},
				Message{
					Topic:     "test",
					Timestamp: time.Now(),
					Payload:   []byte("x"),
				},
				time.NewTimer(time.Hour),
			},
			output: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// ugly, but for this use case specifically we fill the chan
			// It will result in a timeout
			if test.name == "failure_timeout" {
				test.input.sub.messageChan <- Message{}
			}
			res := tryDeliver(test.input.sub, test.input.msg, test.input.timer)
			if res != test.output {
				t.Errorf("res = %v, want %v", res, test.output)
			}
		})
	}
}

func TestUnsubscribeDeleteTopicOnLastSub(t *testing.T) {
	b, _ := New()
	bb := b.(*bus)

	// create a topic and a sub manually
	_ = bb.withWriteState("t", true, func(st *topicState) error {
		st.subs["S"] = &subscription{subscriptionID: "S", messageChan: make(chan Message, 1)}
		return nil
	})

	// build unsubribe fonc and test it
	unsub := bb.buildUnsubscribeFunction("S", "t")
	if err := unsub(); err != nil {
		t.Fatal(err)
	}

	// check if deleted
	bb.mutex.RLock()
	_, ok := bb.subscriptions["t"]
	bb.mutex.RUnlock()
	if ok {
		t.Fatal("topic should be deleted when last sub removed and queue empty")
	}
}

func TestUnsubscribeKeepTopicWhenAutoDeleteDisabled(t *testing.T) {
	b, _ := New(WithAutoDeleteEmptyTopics(false))
	bb := b.(*bus)

	_ = bb.withWriteState("t", true, func(st *topicState) error {
		st.subs["S"] = &subscription{subscriptionID: "S", messageChan: make(chan Message, 1)}
		return nil
	})
	unsub := bb.buildUnsubscribeFunction("S", "t")
	_ = unsub()

	bb.mutex.RLock()
	_, ok := bb.subscriptions["t"]
	bb.mutex.RUnlock()
	if !ok {
		t.Fatal("topic should remain when auto-delete is disabled")
	}
}

func TestRunFanOutDeliveredCounter(t *testing.T) {
	b, _ := New()
	defer b.Close()

	// 1) Abonné réel (buffer 1, drop = true pour que l’envoi ne bloque pas)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sub, err := b.Subscribe(ctx, "t",
		WithBufferSize(1),
		WithDropIfFull(true),
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := b.Publish("t", []byte("x")); err != nil {
		t.Fatal(err)
	}

	// Wait for the read
	select {
	case <-sub.Read():
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting message")
	}

	// ugly but we wait for
	deadline := time.Now().Add(2 * time.Second)
	for {
		st, _ := b.Stats()
		ts, ok := st.PerTopic["t"]
		if ok && ts.Delivered >= 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("expected Delivered>=1, got %d", ts.Delivered)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestRunFanOutDroppedCounter(t *testing.T) {
	b, _ := New()
	bb := b.(*bus)

	_ = bb.withWriteState("t", true, func(st *topicState) error {
		sub := &subscription{
			subscriptionID: "S",
			cfg:            SubscriptionConfig{BufferSize: 1, DropIfFull: true},
			messageChan:    make(chan Message, 1),
		}
		sub.messageChan <- Message{}
		st.subs["S"] = sub
		return nil
	})
	_ = bb.withWriteState("t", true, func(st *topicState) error {
		st.inQueue <- messageRef{topic: "t", payload: []byte("x"), ts: time.Now()}
		return nil
	})
	time.Sleep(20 * time.Millisecond)

	bb.mutex.RLock()
	ts := bb.subscriptions["t"]
	dropped := ts.counters.Dropped.Load()
	bb.mutex.RUnlock()
	if dropped < 1 {
		t.Fatalf("expected Dropped>=1, got %d", dropped)
	}
}
