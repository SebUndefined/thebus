package thebus

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type bus struct {
	mutex     sync.RWMutex
	cfg       *Config
	open      atomic.Bool
	startedAt time.Time
	//subscriptions map[string]*structState
	totals atomicCounters
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
		startedAt: time.Now(),
		cfg:       cfg,
		totals:    atomicCounters{},
		//subscriptions: make(map[string]*structState),
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
