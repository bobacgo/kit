package kafka

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

const (
	addr    = "localhost:9092"
	appName = "example"
	topic   = "test"
)

func TestProducer(t *testing.T) {
	// 创建Kafka服务器配置
	cfg := Config{
		Addrs: []string{addr},
		Producer: ProducerConfig{
			RequiredAcks: RequiredAcksLocal,
			Timeout:      "10s",
			Compression:  CompessionNone,
			Retry: RetryConfig{
				Max:     3,
				Backoff: "1s",
			},
		},
	}

	// 创建Kafka服务器实例
	server := New(appName, cfg)

	// 启动服务器
	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		t.Fatalf("启动Kafka服务器失败: %v", err)
	}

	// 确保在测试结束时关闭服务器
	defer server.Stop(ctx)

	pub, _ := server.Get().(*ProducerServer)

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("Hello, Kafka! %d", i)
		// 发送消息到test主题
		if err := pub.SendMessage(context.Background(), topic, []byte(msg), WithKey(strconv.Itoa(i))); err != nil {
			t.Errorf("发送消息失败: %v", err)
			continue
		}
		time.Sleep(time.Second)
		fmt.Println("消息发送成功！", i)
	}
}

func TestConsumer(t *testing.T) {
	// 创建Kafka服务器配置
	cfg := Config{
		Addrs: []string{addr},
		Consumer: ConsumerConfig{
			OffsetInitial: OffsetNewest,
			MaxWaitTime:   "1s",
			MinBytes:      "1KB",
			MaxBytes:      "10MB",
		},
	}

	// 创建Kafka服务器实例
	server := New(appName, cfg, Subscriber{
		Topic: topic,
		ConsumerInfo: &ConsumerInfo{
			Mode: ModeMutex,
			Handler: func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
				for msg := range claim.Messages() {
					// 打印接收到的消息
					fmt.Printf("收到消息：Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s\n",
						msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))

					// 标记消息已处理
					session.MarkMessage(msg, "")
				}
				return nil
			},
		},
	})

	// 启动服务器
	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		t.Fatalf("启动Kafka服务器失败: %v", err)
	}

	// 确保在测试结束时关闭服务器
	defer server.Stop(ctx)

	// 保持消费者运行
	fmt.Println("消费者已启动，按Ctrl+C退出...")
	select {}
}