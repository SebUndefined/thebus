package thebus

import (
	"context"
)

type Bus interface {
	Publish(topic string, data []byte) (PublishAck, error)
	Subscribe(ctx context.Context, topic string, opts ...SubscribeOption) (Subscription, error)
	Unsubscribe(topic string, subscriberID string) error
	Close() error
	Stats() (StatsResults, error)
}
