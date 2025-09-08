package thebus

import (
	"context"
)

// Bus is the main interface of thebus.
// It provides a publish/subscribe API for sending and receiving messages
// on named topics. A Bus implementation is safe for concurrent use
// by multiple goroutines.
type Bus interface {
	// Publish sends a message on the given topic.
	// If there are subscribers, the message is enqueued in the topic’s buffer
	// and delivered asynchronously. If the buffer is full, behavior depends on
	// configuration: the call could return ErrQueueFull.
	// Returns a PublishAck indicating whether the message was enqueued
	// and how many subscribers were present at publish.
	Publish(topic string, data []byte) (PublishAck, error)
	// Subscribe registers a new subscription to the given topic.
	// A subscription receives all messages published after it is created.
	// Options (buffer size, drop policy, copy strategy, etc.) can be set
	// via SubscribeOption functions.
	// The subscription is automatically unsubscribed when the provided
	// context is canceled.
	Subscribe(ctx context.Context, topic string, opts ...SubscribeOption) (Subscription, error)
	// Unsubscribe removes a subscriber by ID from a topic.
	// It is safe to call multiple times; redundant calls are ignored.
	Unsubscribe(topic string, subscriberID string) error
	// Close shuts down the bus and all topics.
	// After Close, no further publish or subscribe is allowed. Existing
	// subscriptions are closed and workers are stopped. Close is idempotent.
	Close() error
	// Stats returns runtime statistics about the bus, topics, and subscribers.
	Stats() (StatsResults, error)
}
