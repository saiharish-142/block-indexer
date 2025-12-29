package events

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

const (
	TopicEvmBlockCanonical   = "evm.block.canonical"
	TopicEvmTxInserted       = "evm.tx.inserted"
	TopicLogsInserted        = "evm.logs.inserted"
	TopicNewContractDeployed = "contract.deployed"
	TopicDagBlockNew         = "dag.block.new"
	TopicTraceRequested      = "trace.requested"
	TopicTraceCompleted      = "trace.completed"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, payload []byte) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string, handler func(msg []byte) error) (func() error, error)
}

type NATSClient struct {
	conn   *nats.Conn
	prefix string
}

func NewNATSClient(url string, prefix string) (*NATSClient, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect nats: %w", err)
	}
	return &NATSClient{conn: nc, prefix: prefix}, nil
}

func (c *NATSClient) Publish(ctx context.Context, topic string, payload []byte) error {
	subject := c.subject(topic)
	return c.conn.PublishMsg(&nats.Msg{Subject: subject, Data: payload})
}

func (c *NATSClient) Subscribe(ctx context.Context, topic string, handler func(msg []byte) error) (func() error, error) {
	subject := c.subject(topic)
	sub, err := c.conn.Subscribe(subject, func(m *nats.Msg) {
		_ = handler(m.Data)
	})
	if err != nil {
		return nil, fmt.Errorf("subscribe %s: %w", subject, err)
	}
	go func() {
		<-ctx.Done()
		_ = sub.Drain()
	}()
	return sub.Drain, nil
}

func (c *NATSClient) subject(topic string) string {
	if c.prefix == "" {
		return topic
	}
	return fmt.Sprintf("%s.%s", c.prefix, topic)
}
