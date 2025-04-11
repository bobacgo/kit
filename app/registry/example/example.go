package example

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bobacgo/kit/app/registry"
	"github.com/bobacgo/kit/app/registry/consul"
)

// ExampleConsulRegistry 展示如何使用Consul服务注册与发现
func ExampleConsulRegistry() {
	// 创建Consul注册中心客户端
	consulRegistry, err := consul.New(
		consul.WithAddress("127.0.0.1:8500"),
		consul.WithScheme("http"),
		consul.WithTTL(time.Second*15),
		// consul.WithToken("your-token"), // 如果需要认证
	)
	if err != nil {
		log.Fatalf("Failed to create consul registry: %v", err)
	}

	// 创建服务实例
	service := &registry.ServiceInstance{
		Name:      "example-service",
		Version:   "v1.0.0",
		Metadata:  map[string]string{"env": "dev"},
		Endpoints: []string{"http://127.0.0.1:8000", "grpc://127.0.0.1:9000"},
	}

	// 注册服务
	ctx := context.Background()
	if err := consulRegistry.Registry(ctx, service); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	fmt.Printf("Service registered: %s\n", service.String())

	// 获取服务实例
	instances, err := consulRegistry.GetService(ctx, "example-service")
	if err != nil {
		log.Fatalf("Failed to get service: %v", err)
	}

	fmt.Printf("Found %d instances\n", len(instances))
	for _, ins := range instances {
		fmt.Printf("  - %s (version: %s)\n", ins.String(), ins.Version)
	}

	// 监听服务变化
	watcher, err := consulRegistry.Watch(ctx, "example-service")
	if err != nil {
		log.Fatalf("Failed to watch service: %v", err)
	}

	// 在另一个goroutine中处理服务变化
	go func() {
		for {
			instances, err := watcher.Next()
			if err != nil {
				log.Printf("Watcher error: %v", err)
				return
			}

			fmt.Printf("Service instances changed, now have %d instances\n", len(instances))
		}
	}()

	// 等待一段时间后注销服务
	time.Sleep(time.Minute)
	if err := consulRegistry.Deregister(ctx, service); err != nil {
		log.Fatalf("Failed to deregister service: %v", err)
	}

	fmt.Printf("Service deregistered: %s\n", service.String())

	// 停止监听
	watcher.Stop()
}
