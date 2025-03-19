package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
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

// WithKey 影响消息的分区（相同key在同一分区
// 相同 key 的消息进入相同分区，保证消息顺序
// 默认情况下，Kafka 使用 key 的 哈希值 决定消息进入哪个分区
// 如果 key 为空，Kafka 会使用轮询（Round-Robin）方式选择分区
func WithKey(key string) ProducerOpt {
	return func(o *ProducerOpts) {
		o.key = []byte(key)
	}
}

// WithHeaders 元数据(追踪，认证等)
func WithHeaders(headers map[string]string) ProducerOpt {
	return func(o *ProducerOpts) {
		o.headers = headers
	}
}

// WithOffset 影响消息的分区唯一标识 kafka 中的消息
func WithOffset(offset int64) ProducerOpt {
	return func(o *ProducerOpts) {
		o.offset = offset
	}
}

// WithPartition 决定消息存储在那个分区
// Kafka 存储数据的基本单位，每个 Topic 由多个 分区（Partition） 组成
// 影响 Kafka 的并行消费（多个分区可同时被不同消费者处理）
// 保证有序性（同一分区的消息按写入顺序消费）
func WithPartition(partition int32) ProducerOpt {
	return func(o *ProducerOpts) {
		o.partition = partition
	}
}

// WithTimestamp 消息产生时间
// 数据恢复（回溯特定时间点的数据）
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
	if o.offset != 0 {
		msg.Offset = o.offset
	}
	_, _, err := p.pub.SendMessage(msg)
	return err
}
