package kafka

import (
	"github.com/IBM/sarama"
	"github.com/bobacgo/kit/app/types"
)

// Config Kafka服务器配置
type Config struct {
	Addrs []string `mapstructure:"addrs"` // Kafka服务器地址列表

	// 生产者配置
	Producer ProducerConfig `mapstructure:"producer"`
	// 消费者配置
	Consumer ConsumerConfig `mapstructure:"consumer"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	Max     int            `mapstructure:"max" validate:"gte=0"` // 最大重试次数
	Backoff types.Duration `mapstructure:"backoff"`              // 重试间隔
}

type RequiredAcks string

const (
	RequiredAcksAll   RequiredAcks = "all"   // 所有节点确认
	RequiredAcksLocal RequiredAcks = "local" // 本地节点确认
	RequiredAcksNone  RequiredAcks = "none"  // 无需确认
)

func (ack RequiredAcks) ack() sarama.RequiredAcks {
	switch ack {
	case RequiredAcksAll:
		return sarama.WaitForAll
	case RequiredAcksLocal:
		return sarama.WaitForLocal
	case RequiredAcksNone:
		return sarama.NoResponse
	default:
		return sarama.NoResponse
	}
}

type Compersion string

const (
	CompessionNone   Compersion = "none"
	CompessionGzip   Compersion = "gzip"
	CompessionSnappy Compersion = "snappy"
	CompessionLz4    Compersion = "lz4"
)

func (c Compersion) compersion() sarama.CompressionCodec {
	switch c {
	case CompessionGzip:
		return sarama.CompressionGZIP
	case CompessionSnappy:
		return sarama.CompressionSnappy
	case CompessionLz4:
		return sarama.CompressionLZ4
	default:
		return sarama.CompressionNone
	}
}

// ProducerConfig 生产者配置
type ProducerConfig struct {
	// 重试配置
	Retry RetryConfig `mapstructure:"retry"`

	// 发送配置
	RequiredAcks RequiredAcks   `mapstructure:"required_acks" validate:"omitempty,oneof=none local all"`     // 需要的ack数量
	Timeout      types.Duration `mapstructure:"timeout"`                                                     // 发送超时时间
	Compression  Compersion     `mapstructure:"compression" validate:"omitempty,oneof=none gzip snappy lz4"` // 压缩算法
}

type Offset string

const (
	OffsetOldest Offset = "oldest"
	OffsetNewest Offset = "newest"
)

func (o Offset) offset() int64 {
	switch o {
	case OffsetNewest:
		return sarama.OffsetNewest
	case OffsetOldest:
		return sarama.OffsetOldest
	default:
		return sarama.OffsetNewest
	}
}

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	GroupID string `mapstructure:"group_id"` // 消费者组ID

	// 消费配置
	OffsetInitial Offset         `mapstructure:"offset_initial" validate:"omitempty,oneof=oldest newest"` // 初始offset 默认 OffsetNewest
	MaxWaitTime   types.Duration `mapstructure:"max_wait_time"`                                           // 最大等待时间
	MinBytes      types.ByteSize `mapstructure:"min_bytes"`                                               // 最小获取字节数
	MaxBytes      types.ByteSize `mapstructure:"max_bytes"`                                               // 最大获取字节数
}