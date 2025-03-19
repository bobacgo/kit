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
	// 消息可靠性较高，数据不会丢失，延迟较高，吞吐量较低
	RequiredAcksAll RequiredAcks = "all" // 所有节点确认
	// 性能较好，如果 Leader 崩溃，消息可能会丢失
	RequiredAcksLocal RequiredAcks = "local" // Leader节点确认
	// 最高吞吐量，最低延迟，如果 kafka 崩溃，消息可能会丢失
	RequiredAcksNone RequiredAcks = "none" // 无需确认
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

type Compression string

const (
	// 不压缩, 最快，但占用较多带宽， 低延迟应用
	CompessionNone Compression = "none"
	// 压缩率高，但 CPu 开销大 日志/历史数据存储
	CompessionGzip Compression = "gzip"
	// 速度块，CPU占用低，但压缩率一般，实时消息/高吞吐
	CompessionSnappy Compression = "snappy"
	// 比 Snappy 更快，压缩率中等， 高吞吐/低延迟
	CompessionLz4 Compression = "lz4"
)

func (c Compression) compersion() sarama.CompressionCodec {
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
	Compression  Compression    `mapstructure:"compression" validate:"omitempty,oneof=none gzip snappy lz4"` // 压缩算法
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

type GroupRebalance string

const (
	// GroupRebalanceRange 顺序消费：Range（适合处理有序数据，但可能不均衡）
	GroupRebalanceRange GroupRebalance = sarama.RangeBalanceStrategyName
	// GroupRebalanceRoundRobin 负载均衡：RoundRobin（确保消费者负载均衡）
	GroupRebalanceRoundRobin GroupRebalance = sarama.RoundRobinBalanceStrategyName
	// GroupRebalanceSticky 生产环境：Sticky（最稳定，减少 Kafka 重新平衡的影响）
	GroupRebalanceSticky GroupRebalance = sarama.StickyBalanceStrategyName
)

func (r GroupRebalance) rebalance() sarama.BalanceStrategy {
	switch r {
	case GroupRebalanceRange:
		return sarama.NewBalanceStrategyRange()
	case GroupRebalanceRoundRobin:
		return sarama.NewBalanceStrategyRoundRobin()
	case GroupRebalanceSticky:
		return sarama.NewBalanceStrategySticky()
	default:
		return sarama.NewBalanceStrategyRange()
	}
}

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	GroupID        string         `mapstructure:"group_id"`                                                           // 消费者组ID
	GroupRebalance GroupRebalance `mapstructure:"group_rebalance" validate:"omitempty,oneof=range roundrobin sticky"` // 消费者组重新分配策略（当有消费者退出时） 默认 GroupRebalanceRange
	OffsetInitial  Offset         `mapstructure:"offset_initial" validate:"omitempty,oneof=oldest newest"`            // 初始offset 默认 OffsetNewest
	// 最长等待时间，当 MinBytes 还未满足时，最多等待多久再返回数据 默认 250ms
	MaxWaitTime types.Duration `mapstructure:"max_wait_time"` // 最大等待时间
	// 最小拉取字节数，如果消息小于该值，Kafka 可能会等待 默认 1 B
	MinBytes types.ByteSize `mapstructure:"min_bytes"`
	// 最大拉取字节数， 一次 fetch 请求最多拉取多少数据  默认 1MB                                                        // 最小获取字节数
	MaxBytes types.ByteSize `mapstructure:"max_bytes"` // 最大获取字节数
}
