package thebus

import (
	"context"
	"time"
)

type Bus interface {
	Publish(topic string, data []byte) (PublishAck, error)
	Subscribe(ctx context.Context, topic string, opts ...SubscribeOption) (Subscription, error)
	Unsubscribe(topic string, subscriberID string) error
	Close() error
	Stats() (StatsResults, error)
}

// ##############################################################################
// ###############################   PUB/SUB   ##################################
// ##############################################################################

// PublishAck is returned when your client Publish a message on a topic.
// It returns the Topic, if the message is Enqueued or not and the number of
// subscribers. Note that Subscribers is a snapshot when published.
// it may not reflect the real subscribers count
type PublishAck struct {
	Topic       string
	Enqueued    bool
	Subscribers int
}

type Subscription interface {
	GetID() string
	GetTopic() string
	Read() <-chan Message
	Unsubscribe() error
}

type Message struct {
	Topic     string
	Timestamp time.Time
	Payload   []byte
	Seq       uint64
}

// ##############################################################################
// #################################   STATS   ##################################
// ##############################################################################

type Counters struct {
	Published uint64
	Delivered uint64
	Failed    uint64
	Dropped   uint64
}

type StatsResults struct {
	StartedAt   time.Time
	Open        bool
	Topics      int
	Subscribers int
	Totals      Counters
	PerTopic    map[string]TopicStats
}

type TopicStats struct {
	Subscribers int
	Buffered    int
	Counters
}

// ##############################################################################
// #############################   THIRD PARTY   ################################
// ##############################################################################

type Logger interface {
	Info(msg string, kv ...any)
	Error(msg string, kv ...any)
	Debug(msg string, kv ...any)
}

type MetricsHooks interface {
	IncPublished(topic string)
	IncDelivered(topic string)
	IncDropped(topic string)
	IncFailed(topic string)
}
