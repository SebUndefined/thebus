package thebus

import (
	"context"
	"fmt"
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

func New(opts ...Option) (Bus, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	b := &bus{
		startedAt:     time.Now(),
		cfg:           cfg,
		totals:        atomicCounters{},
		subscriptions: make(map[string]*topicState),
	}
	b.open.Store(true)
	return b, nil
}

func (b *bus) Publish(topic string, data []byte) (PublishAck, error) {
	//TODO implement me
	panic("implement me")
}

func (b *bus) Subscribe(ctx context.Context, topic string, opts ...SubscribeOption) (Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (b *bus) Unsubscribe(topic string, subscriberID string) error {
	//TODO implement me
	panic("implement me")
}

func (b *bus) Close() error {
	//TODO implement me
	panic("implement me")
}

func (b *bus) Stats() (StatsResults, error) {
	//TODO implement me
	panic("implement me")
}

// getStateReadLocked
// Caller must hold RLock.
func (b *bus) getStateReadLocked(topic string) (*topicState, bool) {
	state, ok := b.subscriptions[topic]
	return state, ok
}

// getOrCreateStateLocked ensure that the state for a specific topic is initialized.
// If the topic is found, it returns it otherwise, it creates a new one.
// getOrCreateStateLocked start the fanOut method automatically.
// Caller must hold Lock.
func (b *bus) getOrCreateStateLocked(topic string) (*topicState, bool, error) {
	state, ok := b.subscriptions[topic]
	created := false
	if !ok {
		if b.cfg.MaxTopics > 0 && len(b.subscriptions) >= b.cfg.MaxTopics {
			return nil,
				false,
				fmt.Errorf("too many topics (max=%d)", b.cfg.MaxTopics)
		}
		created = true
		// Create the chanel. Apply default in case of bad config
		var queue chan messageRef
		if b.cfg.TopicQueueSize <= 0 {
			queue = make(chan messageRef, DefaultTopicQueueSize)
		} else {
			queue = make(chan messageRef, b.cfg.TopicQueueSize)
		}
		state = &topicState{
			subs:     make(map[string]*subscription),
			counters: atomicCounters{},
			inQueue:  queue,
		}
		b.subscriptions[topic] = state
		if state.started.CompareAndSwap(false, true) {
			state.wg.Add(1)
			// TODO run the fanout like go l.runFanOut(topic, state)
		}
	}
	return state, created, nil
}
