package nom
import "context"
type Event interface {
	GetEventType() string
}
type EventSubscriber interface {
	OnEvent(ctx context.Context, event Event) error
}
