package thebus

type Subscription interface {
	GetID() string
	GetTopic() string
	Read() <-chan Message
	Unsubscribe() error
}

type subscription struct {
	subscriptionID  uuid.UUID
	cfg             SubscriptionConfig
	topic           string
	messageChan     chan Message
	messages        <-chan Message
	unsubscribeFunc func() error
}

var _ Subscription = (*subscription)(nil)

func (s *subscription) GetID() string {
	return s.subscriptionID.String()
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
