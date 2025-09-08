package thebus

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type bus struct {
	mutex         sync.RWMutex
	cfg           *Config
	open          atomic.Bool
	startedAt     time.Time
	subscriptions map[string]*topicState
	totals        atomicCounters
}

var _ Bus = (*bus)(nil)

// New creates a new Bus instance with the given options.
// It applies default values, normalizes the configuration, and starts the bus in an open state.
// The returned Bus is ready to accept Publish and Subscribe calls.
func New(opts ...Option) (Bus, error) {
	cfg := BuildConfig(opts...).Normalize()
	b := &bus{
		startedAt:     time.Now(),
		cfg:           cfg,
		totals:        atomicCounters{},
		subscriptions: make(map[string]*topicState),
	}
	b.open.Store(true)
	return b, nil
}

// Publish a message on a specific topic. It returns a PublishAck and/or an error.
func (b *bus) Publish(topic string, data []byte) (PublishAck, error) {
	now := time.Now().UTC()
	if len(strings.TrimSpace(topic)) == 0 {
		return PublishAck{}, ErrInvalidTopic
	}
	if !b.open.Load() {
		return PublishAck{}, ErrClosed
	}
	var ack PublishAck
	var errOut error
	err := b.withReadState(topic, func(st *topicState) error {
		if st == nil || len(st.subs) == 0 {
			ack = PublishAck{Topic: topic, Enqueued: false, Subscribers: 0}
			return nil
		}
		seq := st.seq.Add(1)

		payload := data
		if b.cfg.CopyOnPublish {
			cp := make([]byte, len(data))
			copy(cp, data)
			payload = cp
		}

		mr := messageRef{
			topic:   topic,
			ts:      now,
			seq:     seq,
			payload: payload,
		}

		select {
		case st.inQueue <- mr:
			st.counters.Published.Add(1)
			b.totals.Published.Add(1)
			ack = PublishAck{
				Topic:       topic,
				Enqueued:    true,
				Subscribers: len(st.subs),
			}
		default:
			ack = PublishAck{
				Topic:       topic,
				Enqueued:    false,
				Subscribers: len(st.subs),
			}
			errOut = ErrQueueFull
		}
		return nil
	})
	if err != nil {
		return PublishAck{}, err
	}
	return ack, errOut
}

func (b *bus) Subscribe(ctx context.Context, topic string, opts ...SubscribeOption) (Subscription, error) {
	// Standard checks
	if !b.open.Load() {
		return nil, ErrClosed
	}
	if len(strings.TrimSpace(topic)) == 0 {
		return nil, ErrInvalidTopic
	}

	// Building the config based on the default one
	// (check in the future the default @SebUndefined)
	cfg := BuildSubscriptionConfig(opts...).Normalize()

	// Building the subscription
	id := b.cfg.IDGenerator()
	msgChan := make(chan Message, cfg.BufferSize)
	sub := &subscription{
		subscriptionID: id,
		cfg:            cfg,
		topic:          topic,
		messageChan:    msgChan,
		messages:       (<-chan Message)(msgChan),
	}
	// Saving, function under lock so ok
	err := b.withWriteState(topic, true, func(state *topicState) error {
		// Recheck in case of closed before the first lock
		// It is possible that someone close it pending we wait for the first lock
		// I do it for avoiding weird state....
		if !b.open.Load() {
			// unlock useless here because it is handled by withWriteState
			return ErrClosed
		}
		if b.cfg.MaxSubscribersPerTopic > 0 && len(state.subs) >= b.cfg.MaxSubscribersPerTopic {
			return fmt.Errorf("too many subscribers per topic (max: %d)", b.cfg.MaxSubscribersPerTopic)
		}
		state.subs[id] = sub
		return nil
	})
	if err != nil {
		return nil, err
	}
	sub.unsubscribeFunc = b.buildUnsubscribeFunction(id, topic)
	go func() {
		select {
		case <-ctx.Done():
			_ = sub.Unsubscribe() // idempotent
		}
	}()

	return sub, nil
}

func (b *bus) buildUnsubscribeFunction(id string, topic string) func() error {
	return func() error {
		b.mutex.Lock()
		if state, ok := b.subscriptions[topic]; ok {
			delete(state.subs, id)
			if len(state.subs) == 0 && len(state.inQueue) == 0 && b.cfg.AutoDeleteEmptyTopics {
				if state.closed.CompareAndSwap(false, true) {
					close(state.inQueue)
					delete(b.subscriptions, topic)
				}
			}
		}
		b.mutex.Unlock()
		return nil
	}
}

func (b *bus) Unsubscribe(topic string, subscriberID string) error {
	// Standard check
	if !b.open.Load() {
		return nil
	}
	if len(strings.TrimSpace(topic)) == 0 {
		return ErrInvalidTopic
	}
	b.mutex.RLock()
	state, ok := b.subscriptions[topic]
	if !ok {
		b.mutex.RUnlock()
		return nil
	}
	sub, ok := state.subs[subscriberID]
	if !ok {
		b.mutex.RUnlock()
		return nil
	}
	b.mutex.RUnlock()
	err := sub.unsubscribeFunc()
	if err != nil {
		return err
	}

	return nil
}

func (b *bus) Close() error {
	// refuse new publish and subscribe
	if !b.open.CompareAndSwap(true, false) {
		return nil
	}
	b.mutex.Lock()
	states := make([]*topicState, 0, len(b.subscriptions))
	for _, st := range b.subscriptions {
		states = append(states, st)
	}
	// close the inQueue of each topic
	for _, st := range states {
		if st.closed.CompareAndSwap(false, true) { // ← idem ici
			close(st.inQueue)
		}
	}
	b.mutex.Unlock()

	// Wait for the worker to stop
	for _, st := range states {
		st.wg.Wait()
	}

	// Cleaning memory
	b.mutex.Lock()
	b.subscriptions = make(map[string]*topicState) // reset propre
	b.mutex.Unlock()

	return nil
}

func (b *bus) Stats() (StatsResults, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	perTopic := make(map[string]TopicStats, len(b.subscriptions))
	subscriberCounts := 0
	for topic, state := range b.subscriptions {
		buffered := 0
		for _, sub := range state.subs {
			subscriberCounts++
			buffered += len(sub.messageChan)
		}
		perTopic[topic] = TopicStats{
			Subscribers: len(state.subs),
			Buffered:    buffered,
			Counters: Counters{
				Published: state.counters.Published.Load(),
				Delivered: state.counters.Delivered.Load(),
				Failed:    state.counters.Failed.Load(),
				Dropped:   state.counters.Dropped.Load(),
			},
		}
	}
	s := StatsResults{
		StartedAt:   b.startedAt,
		Open:        b.open.Load(),
		Topics:      len(b.subscriptions),
		Subscribers: subscriberCounts,
		Totals: Counters{
			Published: b.totals.Published.Load(),
			Delivered: b.totals.Delivered.Load(),
			Failed:    b.totals.Failed.Load(),
			Dropped:   b.totals.Dropped.Load(),
		},
		PerTopic: perTopic,
	}
	return s, nil
}

func (b *bus) withReadState(topic string, callback func(st *topicState) error) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return callback(b.subscriptions[topic])
}

func (b *bus) withWriteState(topic string, createIfNotExists bool, writeFunc func(state *topicState) error) error {
	b.mutex.Lock()
	state, ok := b.subscriptions[topic]
	start := false
	if !ok {
		if !createIfNotExists {
			err := writeFunc(nil)
			b.mutex.Unlock()
			return err
		}
		if b.cfg.MaxTopics > 0 && len(b.subscriptions) >= b.cfg.MaxTopics {
			b.mutex.Unlock()
			return fmt.Errorf("too many topics (max=%d)", b.cfg.MaxTopics)
		}
		qSize := b.cfg.TopicQueueSize
		if qSize <= 0 {
			qSize = DefaultTopicQueueSize
		}
		state = newTopicState(qSize)
		// add the state
		b.subscriptions[topic] = state

		if state.started.CompareAndSwap(false, true) {
			state.wg.Add(1)
			start = true
		}
	}
	err := writeFunc(state)
	b.mutex.Unlock()
	if start {
		go b.runFanOut(topic, state)
	}
	return err
}
