package consul

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bobacgo/kit/app/registry"
	"github.com/bobacgo/kit/pkg/uid"
	"github.com/hashicorp/consul/api"
)

var (
	ErrServiceInstanceNotFound = errors.New("service instance not found")
	ErrWatcherStopped          = errors.New("watcher stopped")
)

// Option is consul registry option
type Option func(o *options)

// options is consul registry options
type options struct {
	address string
	scheme  string
	token   string
	ttl     time.Duration
}

// WithAddress with registry address
func WithAddress(address string) Option {
	return func(o *options) {
		o.address = address
	}
}

// WithScheme with registry scheme
func WithScheme(scheme string) Option {
	return func(o *options) {
		o.scheme = scheme
	}
}

// WithToken with registry token
func WithToken(token string) Option {
	return func(o *options) {
		o.token = token
	}
}

// WithTTL with registry ttl
func WithTTL(ttl time.Duration) Option {
	return func(o *options) {
		o.ttl = ttl
	}
}

// Registry is consul registry
type Registry struct {
	opt       options
	client    *api.Client
	instances map[string]*registry.ServiceInstance
	sync.RWMutex
}

// New create a consul registry
func New(opts ...Option) (*Registry, error) {
	options := options{
		address: "127.0.0.1:8500",
		scheme:  "http",
		ttl:     time.Second * 15,
	}
	for _, o := range opts {
		o(&options)
	}

	config := api.DefaultConfig()
	config.Address = options.address
	config.Scheme = options.scheme
	if options.token != "" {
		config.Token = options.token
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Registry{
		opt:       options,
		client:    client,
		instances: make(map[string]*registry.ServiceInstance),
	}, nil
}

// Registry 注册服务
func (r *Registry) Registry(ctx context.Context, service *registry.ServiceInstance) error {
	if service.ID == "" {
		service.ID = uid.UUID()
	}

	// 创建Consul服务注册信息
	registration := &api.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Tags:    []string{service.Version},
		Address: extractAddress(service.Endpoints),
		Port:    extractPort(service.Endpoints),
		Meta:    service.Metadata,
	}

	// 添加健康检查
	check := &api.AgentServiceCheck{
		TTL:                            r.opt.ttl.String(),
		DeregisterCriticalServiceAfter: "30s",
	}
	registration.Check = check

	// 注册服务
	if err := r.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	// 保存服务实例到本地缓存
	r.Lock()
	r.instances[service.ID] = service
	r.Unlock()

	// 启动TTL健康检查更新
	go r.ttlHealthCheck(service)

	return nil
}

// ttlHealthCheck 定期更新服务健康状态
func (r *Registry) ttlHealthCheck(service *registry.ServiceInstance) {
	ticker := time.NewTicker(r.opt.ttl / 2)
	defer ticker.Stop()

	for {
		<-ticker.C
		r.RLock()
		_, ok := r.instances[service.ID]
		r.RUnlock()
		if !ok {
			return
		}

		err := r.client.Agent().UpdateTTL("service:"+service.ID, "healthy", api.HealthPassing)
		if err != nil {
			// 如果更新失败，尝试重新注册
			r.Registry(context.Background(), service)
		}
	}
}

// Deregister 注销服务
func (r *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	err := r.client.Agent().ServiceDeregister(service.ID)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %v", err)
	}

	// 从本地缓存中移除服务实例
	r.Lock()
	delete(r.instances, service.ID)
	r.Unlock()

	return nil
}

// GetService 获取服务实例
func (r *Registry) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	// 查询服务实例
	services, _, err := r.client.Health().Service(serviceName, "", true, (&api.QueryOptions{}).WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %v", err)
	}

	if len(services) == 0 {
		return nil, nil
	}

	instances := make([]*registry.ServiceInstance, 0, len(services))
	for _, service := range services {
		var version string
		if len(service.Service.Tags) > 0 {
			version = service.Service.Tags[0]
		}

		endpoint := fmt.Sprintf("http://%s:%d", service.Service.Address, service.Service.Port)
		instance := &registry.ServiceInstance{
			ID:        service.Service.ID,
			Name:      service.Service.Service,
			Version:   version,
			Metadata:  service.Service.Meta,
			Endpoints: []string{endpoint},
		}

		instances = append(instances, instance)
	}

	return instances, nil
}

// Watch 根据服务名称创建观察者
func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return newWatcher(ctx, r, serviceName), nil
}

// consulWatcher 实现 registry.Watcher 接口
type consulWatcher struct {
	registry    *Registry
	serviceName string
	watchCh     chan []*registry.ServiceInstance
	stopCh      chan struct{}
	stopped     bool
	mutex       sync.Mutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// newWatcher 创建一个新的consul观察者
func newWatcher(ctx context.Context, r *Registry, serviceName string) *consulWatcher {
	ctx, cancel := context.WithCancel(ctx)
	w := &consulWatcher{
		registry:    r,
		serviceName: serviceName,
		watchCh:     make(chan []*registry.ServiceInstance, 100),
		stopCh:      make(chan struct{}),
		ctx:         ctx,
		cancel:      cancel,
	}

	go w.watch()
	return w
}

// watch 监听服务变化
func (w *consulWatcher) watch() {
	var lastIndex uint64 = 0
	for {
		select {
		case <-w.stopCh:
			return
		default:
			// 使用阻塞查询监听服务变化
			services, meta, err := w.registry.client.Health().Service(
				w.serviceName,
				"",
				true,
				(&api.QueryOptions{
					WaitIndex: lastIndex,
				}).WithContext(w.ctx),
			)

			if err != nil {
				if w.ctx.Err() != nil {
					return
				}
				time.Sleep(time.Second)
				continue
			}

			// 检查索引是否变化
			if lastIndex != meta.LastIndex {
				lastIndex = meta.LastIndex

				instances := make([]*registry.ServiceInstance, 0, len(services))
				for _, service := range services {
					var version string
					if len(service.Service.Tags) > 0 {
						version = service.Service.Tags[0]
					}

					endpoint := fmt.Sprintf("http://%s:%d", service.Service.Address, service.Service.Port)
					instance := &registry.ServiceInstance{
						ID:        service.Service.ID,
						Name:      service.Service.Service,
						Version:   version,
						Metadata:  service.Service.Meta,
						Endpoints: []string{endpoint},
					}

					instances = append(instances, instance)
				}

				select {
				case w.watchCh <- instances:
				case <-w.stopCh:
					return
				case <-w.ctx.Done():
					return
				}
			}
		}
	}
}

// Next 监听服务实例变化
func (w *consulWatcher) Next() ([]*registry.ServiceInstance, error) {
	w.mutex.Lock()
	if w.stopped {
		w.mutex.Unlock()
		return nil, ErrWatcherStopped
	}
	w.mutex.Unlock()

	select {
	case instances := <-w.watchCh:
		return instances, nil
	case <-w.stopCh:
		return nil, ErrWatcherStopped
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	}
}

// Stop 停止监听
func (w *consulWatcher) Stop() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.stopped {
		w.stopped = true
		w.cancel()
		close(w.stopCh)
	}

	return nil
}

// 辅助函数：从endpoints中提取地址
func extractAddress(endpoints []string) string {
	if len(endpoints) == 0 {
		return ""
	}

	// 简单实现，实际应用中应该更加健壮
	// 假设格式为 http://127.0.0.1:8000 或 grpc://127.0.0.1:9000
	for _, endpoint := range endpoints {
		// 尝试解析地址和端口
		host := ""
		if len(endpoint) > 7 && endpoint[:7] == "http://" {
			host = endpoint[7:]
		} else if len(endpoint) > 7 && endpoint[:7] == "grpc://" {
			host = endpoint[7:]
		}

		if host != "" {
			// 分离地址和端口
			for i := 0; i < len(host); i++ {
				if host[i] == ':' {
					return host[:i]
				}
			}
		}
	}

	return ""
}

// 辅助函数：从endpoints中提取端口
func extractPort(endpoints []string) int {
	if len(endpoints) == 0 {
		return 0
	}

	// 简单实现，实际应用中应该更加健壮
	// 假设格式为 http://127.0.0.1:8000 或 grpc://127.0.0.1:9000
	for _, endpoint := range endpoints {
		// 尝试解析地址和端口
		host := ""
		if len(endpoint) > 7 && endpoint[:7] == "http://" {
			host = endpoint[7:]
		} else if len(endpoint) > 7 && endpoint[:7] == "grpc://" {
			host = endpoint[7:]
		}

		if host != "" {
			// 分离地址和端口
			for i := 0; i < len(host); i++ {
				if host[i] == ':' {
					port := 0
					fmt.Sscanf(host[i+1:], "%d", &port)
					return port
				}
			}
		}
	}

	return 0
}
