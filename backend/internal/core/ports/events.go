package ports

type EventPublisher interface {
	PublishEvent(channel string, event interface{}) error
}
