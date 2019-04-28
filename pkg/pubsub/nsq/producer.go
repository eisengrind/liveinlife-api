package nsq

import (
	"context"

	"github.com/51st-state/api/pkg/pubsub"
	"github.com/nsqio/go-nsq"
)

type producer struct {
	nsq   *nsq.Producer
	topic string
}

// NewProducer for the nsq message queue
func NewProducer(n *nsq.Producer, topic string) pubsub.Producer {
	return &producer{
		nsq:   n,
		topic: topic,
	}
}

func (p *producer) Produce(ctx context.Context, data []byte) error {
	return p.nsq.Publish(p.topic, data)
}
