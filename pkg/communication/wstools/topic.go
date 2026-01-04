package wstools

import (
	"context"
	"log/slog"
	"strings"

	"github.com/coder/websocket"
	"github.com/xdoubleu/essentia/pkg/threading"
)

// OnSubscribeCallback is called to fetch data that
// should be returned when a new subscriber is added to a topic.
type OnSubscribeCallback = func(ctx context.Context, topic *Topic) (any, error)

// Topic is used to efficiently send messages
// to [Subscriber]s in a WebSocket.
type Topic struct {
	Name                string
	allowedOrigins      []string
	eventQueue          *threading.EventQueue
	onSubscribeCallback OnSubscribeCallback
}

// NewTopic creates a new [Topic].
func NewTopic(
	logger *slog.Logger,
	name string,
	allowedOrigins []string,
	maxWorkers int,
	channelBufferSize int,
	onSubscribeCallback OnSubscribeCallback,
) *Topic {
	for i, url := range allowedOrigins {
		if strings.Contains(url, "://") {
			allowedOrigins[i] = strings.Split(url, "://")[1]
		}
	}

	return &Topic{
		Name:           name,
		allowedOrigins: allowedOrigins,
		eventQueue: threading.NewEventQueue(
			logger,
			maxWorkers,
			channelBufferSize,
		),
		onSubscribeCallback: onSubscribeCallback,
	}
}

// Subscribe subscribes a [Subscriber] to this [Topic].
// If configured a message will be sent on subscribing.
// If no message handling go routine was
// running this will be started now.
func (t *Topic) Subscribe(conn *websocket.Conn) error {
	sub := NewSubscriber(t, conn)
	t.eventQueue.AddSubscriber(sub)

	if t.onSubscribeCallback != nil {
		event, err := t.onSubscribeCallback(context.Background(), t)
		if err != nil {
			return err
		}

		sub.OnEventCallback(event)
	}

	return nil
}

// UnSubscribe unsubscribes a [Subscriber] from this [Topic].
func (t *Topic) UnSubscribe(sub Subscriber) {
	t.eventQueue.RemoveSubscriber(sub)
}

// EnqueueEvent enqueues an event if there are subscribers on this [Topic].
func (t *Topic) EnqueueEvent(event any) {
	t.eventQueue.EnqueueEvent(event)
}
