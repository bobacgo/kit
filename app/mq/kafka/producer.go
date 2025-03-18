package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"time"
)

type ProducerOpt func(o *ProducerOpts)

type ProducerOpts struct {
	key       []byte
	headers   map[string]string
	offset    int64
	partition int32
	timestamp time.Time
}

func (o ProducerOpts) h() []sarama.RecordHeader {
	headers := make([]sarama.RecordHeader, 0, len(o.headers))
	for k, v := range o.headers {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}
	return headers
}

func WithKey(key string) ProducerOpt {
	return func(o *ProducerOpts) {
		o.key = []byte(key)
	}
}
func WithHeaders(headers map[string]string) ProducerOpt {
	return func(o *ProducerOpts) {
		o.headers = headers
	}
}
func WithOffset(offset int64) ProducerOpt {
	return func(o *ProducerOpts) {
		o.offset = offset
	}
}
func WithPartition(partition int32) ProducerOpt {
	return func(o *ProducerOpts) {
		o.partition = partition
	}
}
func WithTimestamp(timestamp time.Time) ProducerOpt {
	return func(o *ProducerOpts) {
		o.timestamp = timestamp
	}
}

type ProducerServer struct {
	pub sarama.SyncProducer
}

func newProducer(adds []string, cfg *ProducerConfig) (*ProducerServer, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = cfg.RequiredAcks.ack()
	if cfg.Timeout != "" {
		config.Producer.Timeout = cfg.Timeout.TimeDuration()
	}
	if cfg.Retry.Max > 0 {
		config.Producer.Retry.Max = cfg.Retry.Max
	}
	if cfg.Retry.Backoff != "" {
		config.Producer.Retry.Backoff = cfg.Retry.Backoff.TimeDuration()
	}
	if cfg.Compression != "" {
		config.Producer.Compression = cfg.Compression.compersion()
	}

	// 确保消息发送成功
	config.Producer.Return.Successes = true
	pub, err := sarama.NewSyncProducer(adds, config)
	if err != nil {
		return nil, err
	}

	return &ProducerServer{
		pub: pub,
	}, nil
}

// SendMessage 发送消息到指定主题
func (p *ProducerServer) SendMessage(ctx context.Context, topic string, value []byte, opts ...ProducerOpt) error {

	// TODO ctx to headers

	o := ProducerOpts{}
	for _, opt := range opts {
		opt(&o)
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	if len(o.key) > 0 {
		msg.Key = sarama.ByteEncoder(o.key)
	}
	if o.partition != 0 {
		msg.Partition = o.partition
	}
	if !o.timestamp.IsZero() {
		msg.Timestamp = o.timestamp
	}
	if len(o.headers) > 0 {
		msg.Headers = o.h()
	}
	_, _, err := p.pub.SendMessage(msg)
	return err
}