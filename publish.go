package thebus

// PublishAck is returned when your client Publish a message on a topic.
// It returns the Topic, if the message is Enqueued or not and the number of
// subscribers. Note that Subscribers is a snapshot when published.
// it may not reflect the real subscribers count
type PublishAck struct {
	Topic       string
	Enqueued    bool
	Subscribers int
}
