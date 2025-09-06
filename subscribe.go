package thebus

import (
	"sync"
	"sync/atomic"
	"time"
)

type Subscription interface {
	GetID() string
	GetTopic() string
	Read() <-chan Message
	Unsubscribe() error
}

type SubscriptionConfig struct {
	Strategy    SubscriptionStrategy
	BufferSize  int
	SendTimeout time.Duration
	DropIfFull  bool
}

type subscription struct {
	subscriptionID  string
	cfg             SubscriptionConfig
	topic           string
	messageChan     chan Message
	messages        <-chan Message
	unsubscribeFunc func() error
}

func DefaultSubscriptionConfig() SubscriptionConfig {
	return SubscriptionConfig{
		BufferSize:  128,
		SendTimeout: 200 * time.Millisecond,
		DropIfFull:  true,
		Strategy:    SubscriptionStrategyPayloadShared,
	}
}

var _ Subscription = (*subscription)(nil)

func (s *subscription) GetID() string {
	return s.subscriptionID
}

func (s *subscription) GetTopic() string {
	return s.topic
}

func (s *subscription) Read() <-chan Message {
	return s.messages
}

func (s *subscription) Unsubscribe() error {
	return s.unsubscribeFunc()
}

// SubscribeOption -
type SubscribeOption func(subCfg *SubscriptionConfig)

func WithStrategy(strategy SubscriptionStrategy) SubscribeOption {
	return func(subCfg *SubscriptionConfig) {
		subCfg.Strategy = strategy
	}
}

func WithBufferSize(bufferSize int) SubscribeOption {
	return func(subCfg *SubscriptionConfig) {
		if bufferSize < 1 {
			return
		}
		subCfg.BufferSize = bufferSize
	}
}
func WithSendTimeout(timeout time.Duration) SubscribeOption {
	return func(subCfg *SubscriptionConfig) {
		subCfg.SendTimeout = timeout
	}
}

func WithDropIfFull(dropIfFull bool) SubscribeOption {
	return func(subCfg *SubscriptionConfig) {
		subCfg.DropIfFull = dropIfFull
	}
}

type topicState struct {
	subs     map[string]*subscription
	started  atomic.Bool
	counters atomicCounters
	inQueue  chan messageRef
	seq      atomic.Uint64
	wg       sync.WaitGroup
	closed   atomic.Bool
}
