package thebus

import "errors"

var (
	ErrClosed                   = errors.New("thebus.closed")
	ErrQueueFull                = errors.New("thebus.queue.full")
	ErrInvalidTopic             = errors.New("thebus.invalid.topic")
	ErrInvalidSubscriberID      = errors.New("thebus.invalid.subscriberID")
	ErrInvalidTopicName         = errors.New("thebus.invalid.topic.name")
	ErrInvalidTopicNameReserved = errors.New("thebus.invalid.topic.name.reserved")
	ErrIDGeneratorNotSet        = errors.New("thebus.idGenerator.not_set")
)
