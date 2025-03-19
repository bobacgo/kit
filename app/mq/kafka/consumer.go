package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/bobacgo/kit/pkg/uid"
)

type Subscriber struct {
	*ConsumerInfo
	Topic string
}

// ConsumerMode 消费模式
type ConsumerMode int

const (
	// ModeMutex 互斥模式：消息只被一个消费者消费
	ModeMutex ConsumerMode = iota
	// ModeBroadcast 广播模式：所有消费者都能收到消息
	ModeBroadcast
)

// ConsumerInfo 消费者处理器信息
type ConsumerInfo struct {
	Handler ConsumeHandlerFunc
	Mode    ConsumerMode
}

type ConsumeHandlerFunc func(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error

type consumerServer struct {
	subs     map[string]sarama.ConsumerGroup // 每个主题一个消费者组
	handlers map[string]*ConsumerInfo        // 每个主题的处理器信息
}

func newConsumer(addrs []string, cfg ConsumerConfig, handlers map[string]*ConsumerInfo) (*consumerServer, error) {
	config := sarama.NewConfig()

	if cfg.GroupRebalance != "" {
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{cfg.GroupRebalance.rebalance()}
	}
	if cfg.OffsetInitial != "" {
		config.Consumer.Offsets.Initial = cfg.OffsetInitial.offset()
	}
	if cfg.MaxWaitTime != "" {
		config.Consumer.MaxWaitTime = cfg.MaxWaitTime.TimeDuration()
	}
	if cfg.MinBytes != "" {
		config.Consumer.Fetch.Min = int32(cfg.MinBytes.Int())
	}
	if cfg.MaxBytes != "" {
		config.Consumer.Fetch.Max = int32(cfg.MaxBytes.Int())
	}

	c := &consumerServer{
		handlers: handlers,
		subs:     make(map[string]sarama.ConsumerGroup),
	}

	// 按消费模式分组主题
	mutexTopics := make(map[string][]string)   // groupID -> topics
	broadcastTopics := make(map[string]string) // topic -> groupID

	// 为每个主题准备消费者组ID
	for topic, info := range handlers {
		groupID := cfg.GroupID
		if info.Mode == ModeBroadcast {
			groupID = fmt.Sprintf("%s-%s-%s", groupID, topic, uid.UUID())
			broadcastTopics[topic] = groupID
		} else {
			groupID = fmt.Sprintf("%s-mutex", groupID)
			mutexTopics[groupID] = append(mutexTopics[groupID], topic)
		}
	}

	// 创建互斥模式的消费者组
	if len(mutexTopics) > 0 {
		// 只需创建一个消费者组连接
		consumer, err := sarama.NewConsumerGroup(addrs, fmt.Sprintf("%s-mutex", cfg.GroupID), config)
		if err != nil {
			return nil, fmt.Errorf("创建互斥模式消费者组失败: %w", err)
		}

		// 收集所有互斥模式的主题
		var allTopics []string
		for _, topics := range mutexTopics {
			allTopics = append(allTopics, topics...)
			// 建立主题到消费者的映射
			for _, topic := range topics {
				c.subs[topic] = consumer
			}
		}
		// 启动一个消费循环处理所有互斥主题
		go c.consumeLoop(allTopics...)
	}

	// 创建广播模式的消费者组
	for topic, groupID := range broadcastTopics {
		consumer, err := sarama.NewConsumerGroup(addrs, groupID, config)
		if err != nil {
			return nil, fmt.Errorf("创建广播模式消费者组失败: %w", err)
		}
		c.subs[topic] = consumer
		go c.consumeLoop(topic)
	}

	return c, nil
}

func (srv *consumerServer) consumeLoop(topics ...string) {
	for {
		// 获取消费者组（使用第一个主题作为key）
		consumer, ok := srv.subs[topics[0]]

		if !ok {
			time.Sleep(time.Second)
			continue
		}

		// 消费消息
		if err := consumer.Consume(context.Background(), topics, srv); err != nil {
			slog.Error("消费消息出错", "topics", topics, "error", err)
		}
	}
}

// 实现 sarama.ConsumerGroupHandler 接口
func (srv *consumerServer) Setup(session sarama.ConsumerGroupSession) error {
	slog.Info("正在监听:", "topics", session.Claims())
	return nil
}

func (srv *consumerServer) Cleanup(session sarama.ConsumerGroupSession) error {
	slog.Info("取消监听:", "topics", session.Claims(), "member_id", session.MemberID())
	return nil
}

func (srv *consumerServer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	info, ok := srv.handlers[claim.Topic()]
	if !ok {
		return fmt.Errorf("未找到主题 %s 的处理器", claim.Topic())
	}
	return info.Handler(session, claim)
}
