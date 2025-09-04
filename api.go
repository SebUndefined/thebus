package thebus

import "context"

type Bus interface {
	Publish(topic string, data []byte) (PublishAck, error)
	Subscribe(ctx context.Context, topic string, opts ...SubscribeOption) (Subscription, error)
	Unsubscribe(topic string, subscriberID string) error
	Close() error
	Stats() (StatsResults, error)
}

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
