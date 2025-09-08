package thebus

type MetricsHooks interface {
	IncPublished(topic string)
	IncDelivered(topic string)
	IncDropped(topic string)
	IncFailed(topic string)
}
