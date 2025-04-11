package registry

import "context"

// ServiceRegistrar 服务注册
type ServiceRegistrar interface {
	// Registry 注册服务
	Registry(ctx context.Context, service *ServiceInstance) error
	// Deregister 注销服务
	Deregister(ctx context.Context, service *ServiceInstance) error
}

// ServiceDiscovery 服务发现
// 1.本地缓存 (不需要每次请求服务,都去注册中心拿取)
// 2.与注册中心长连接
// 3.服务实例发生变化,直接推送给订阅端.
type ServiceDiscovery interface {
	// GetService 获取服务实例
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	// Watch 根据服务名称创建观察者
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

type Watcher interface {
	// Next 监听服务实例变化
	Next() ([]*ServiceInstance, error)
	// Stop 停止监听
	Stop() error
}

type ServiceInstance struct {
	ID       string            `json:"id"`       // 服务ID
	Name     string            `json:"name"`     // 服务名称
	Version  string            `json:"version"`  // 服务版本
	Metadata map[string]string `json:"metadata"` // 服务元数据

	// http://127.0.0.1:8000
	// grpc://127.0.0.1:9000
	Endpoints []string `json:"endpoints"`
}

func (s *ServiceInstance) String() string {
	return s.Name + "-" + s.ID
}
