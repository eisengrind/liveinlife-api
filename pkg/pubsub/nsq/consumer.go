package nsq

import (
	"context"

	"github.com/51st-state/api/pkg/pubsub"
	"github.com/nsqio/go-nsq"
)

type consumer struct {
	nsq        *nsq.Consumer
	lookupAddr string
}

// NewConsumer on a nsq message queue
func NewConsumer(topic, channel, lookupAddr string, cfg *nsq.Config) (pubsub.Consumer, error) {
	n, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		return nil, err
	}

	return &consumer{
		n,
		lookupAddr,
	}, nil
}

type message struct {
	*nsq.Message
}

func newMessage(msg *nsq.Message) pubsub.Message {
	return &message{
		msg,
	}
}

func (m *message) Ack() {
	m.Finish()
}

func (m *message) Nack() {
	m.Requeue(-1)
}

func (m *message) Data() []byte {
	return m.Body
}

type consumerHandler struct {
	handlerFunc pubsub.HandlerFunc
}

func newConsumerHandler(h pubsub.HandlerFunc) nsq.Handler {
	return &consumerHandler{
		handlerFunc: h,
	}
}

func (c *consumerHandler) HandleMessage(message *nsq.Message) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return c.handlerFunc(ctx, newMessage(message))
}

func (c *consumer) Consume(ctx context.Context, h pubsub.HandlerFunc) error {
	c.nsq.AddHandler(newConsumerHandler(h))
	if err := c.nsq.ConnectToNSQLookupd(c.lookupAddr); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
