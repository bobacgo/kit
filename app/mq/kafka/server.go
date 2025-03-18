package kafka

import (
	"context"
	"fmt"
	"sync"
)

// Server 实现 server.Server 接口的 Kafka 服务器
type Server struct {
	appName string
	config  *Config

	// 生产者
	producer *ProducerServer

	// 消费者
	consMu     sync.RWMutex
	consumer   *consumerServer
	handlers   map[string]*ConsumerInfo // 每个主题的处理器信息
	consClosed chan struct{}
}

// New 创建新的Kafka服务器实例
func New(appName string, config Config, subs ...Subscriber) *Server {
	handlers := make(map[string]*ConsumerInfo)
	for _, sub := range subs {
		handlers[sub.Topic] = sub.ConsumerInfo
	}
	return &Server{
		appName:    appName,
		config:     &config,
		handlers:   handlers,
		consClosed: make(chan struct{}),
	}
}

// Start 启动Kafka服务器
func (s *Server) Start(ctx context.Context) error {
	// 初始化生产者
	var err error
	if s.producer, err = newProducer(s.config.Addrs, &s.config.Producer); err != nil {
		return fmt.Errorf("init producer : %w", err)
	}

	if s.config.Consumer.GroupID == "" {
		s.config.Consumer.GroupID = s.appName
	}

	// 初始化消费者
	if s.consumer, err = newConsumer(s.config.Addrs, s.config.Consumer, s.handlers); err != nil {
		return fmt.Errorf("init consumer: %w", err)
	}
	return nil
}

// Stop 停止Kafka服务器
func (s *Server) Stop(ctx context.Context) error {
	// 关闭生产者
	if s.producer != nil {
		s.producer.pub.Close()
	}

	// 关闭所有消费者
	s.consMu.Lock()
	for _, consumer := range s.consumer.subs {
		if consumer != nil {
			consumer.Close()
		}
	}
	s.consMu.Unlock()

	// 等待消费者完全关闭
	<-s.consClosed
	return nil
}

// Get 获取组件实例
func (s *Server) Get() any {
	return s.producer
}
