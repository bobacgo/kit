package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bobacgo/kit/app/registry"
	"github.com/bobacgo/kit/pkg/uid"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	ErrServiceInstanceNotFound = errors.New("service instance not found")
	ErrWatcherStopped          = errors.New("watcher stopped")
)

// Option is etcd registry option
type Option func(o *options)

// options is etcd registry options
type options struct {
	endpoints []string
	username  string
	password  string
	timeout   time.Duration
	ttl       time.Duration
}

// WithEndpoints with registry endpoints
func WithEndpoints(endpoints ...string) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

// WithUsername with registry username
func WithUsername(username string) Option {
	return func(o *options) {
		o.username = username
	}
}

// WithPassword with registry password
func WithPassword(password string) Option {
	return func(o *options) {
		o.password = password
	}
}

// WithTimeout with registry timeout
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithTTL with registry ttl
func WithTTL(ttl time.Duration) Option {
	return func(o *options) {
		o.ttl = ttl
	}
}

// Registry is etcd registry
type Registry struct {
	opt       options
	client    *clientv3.Client
	instances map[string]*registry.ServiceInstance
	lease     clientv3.Lease
	leaseID   clientv3.LeaseID
	sync.RWMutex
}

// New create a etcd registry
func New(opts ...Option) (*Registry, error) {
	options := options{
		endpoints: []string{"127.0.0.1:2379"},
		timeout:   time.Second * 5,
		ttl:       time.Second * 15,
	}
	for _, o := range opts {
		o(&options)
	}

	config := clientv3.Config{
		Endpoints:   options.endpoints,
		DialTimeout: options.timeout,
	}
	if options.username != "" && options.password != "" {
		config.Username = options.username
		config.Password = options.password
	}

	client, err := clientv3.New(config)
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

	// 创建租约
	lease := clientv3.NewLease(r.client)
	leaseResp, err := lease.Grant(ctx, int64(r.opt.ttl.Seconds()))
	if err != nil {
		return fmt.Errorf("failed to create lease: %v", err)
	}

	// 保存租约ID
	r.lease = lease
	r.leaseID = leaseResp.ID

	// 序列化服务实例
	data, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service: %v", err)
	}

	// 构建服务键
	key := fmt.Sprintf("/services/%s/%s", service.Name, service.ID)

	// 注册服务
	_, err = r.client.Put(ctx, key, string(data), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	// 保存服务实例到本地缓存
	r.Lock()
	r.instances[service.ID] = service
	r.Unlock()

	// 启动自动续约
	go r.keepAlive(ctx, service)

	return nil
}

// keepAlive 自动续约
func (r *Registry) keepAlive(ctx context.Context, service *registry.ServiceInstance) {
	// 创建租约keepalive
	keepaliveCh, err := r.lease.KeepAlive(ctx, r.leaseID)
	if err != nil {
		// 如果创建失败，尝试重新注册
		r.Registry(context.Background(), service)
		return
	}

	for {
		select {
		case _, ok := <-keepaliveCh:
			if !ok {
				// 如果通道关闭，尝试重新注册
				r.RLock()
				_, exists := r.instances[service.ID]
				r.RUnlock()
				if exists {
					r.Registry(context.Background(), service)
				}
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// Deregister 注销服务
func (r *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	// 构建服务键
	key := fmt.Sprintf("/services/%s/%s", service.Name, service.ID)

	// 删除服务
	_, err := r.client.Delete(ctx, key)
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
	// 构建服务前缀
	prefix := fmt.Sprintf("/services/%s/", serviceName)

	// 查询服务实例
	resp, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %v", err)
	}

	if len(resp.Kvs) == 0 {
		return nil, nil
	}

	instances := make([]*registry.ServiceInstance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		instance := &registry.ServiceInstance{}
		if err := json.Unmarshal(kv.Value, instance); err != nil {
			return nil, fmt.Errorf("failed to unmarshal service: %v", err)
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// Watch 根据服务名称创建观察者
func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return newWatcher(ctx, r, serviceName), nil
}

// etcdWatcher 实现 registry.Watcher 接口
type etcdWatcher struct {
	registry    *Registry
	serviceName string
	watchCh     chan []*registry.ServiceInstance
	stopCh      chan struct{}
	stopped     bool
	mutex       sync.Mutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// newWatcher 创建一个新的etcd观察者
func newWatcher(ctx context.Context, r *Registry, serviceName string) *etcdWatcher {
	ctx, cancel := context.WithCancel(ctx)
	w := &etcdWatcher{
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
func (w *etcdWatcher) watch() {
	// 构建服务前缀
	prefix := fmt.Sprintf("/services/%s/", w.serviceName)

	// 获取当前服务列表
	resp, err := w.registry.client.Get(w.ctx, prefix, clientv3.WithPrefix())
	if err == nil && resp.Count > 0 {
		instances := make([]*registry.ServiceInstance, 0, len(resp.Kvs))
		for _, kv := range resp.Kvs {
			instance := &registry.ServiceInstance{}
			if err := json.Unmarshal(kv.Value, instance); err == nil {
				instances = append(instances, instance)
			}
		}

		select {
		case w.watchCh <- instances:
		case <-w.stopCh:
			return
		case <-w.ctx.Done():
			return
		}
	}

	// 创建watcher
	watchCh := w.registry.client.Watch(w.ctx, prefix, clientv3.WithPrefix(), clientv3.WithRev(resp.Header.Revision+1))
	for {
		select {
		case <-w.stopCh:
			return
		case <-w.ctx.Done():
			return
		case watchResp := <-watchCh:
			if watchResp.Canceled {
				return
			}

			for range watchResp.Events {
				// 服务发生变化，重新获取服务列表
				resp, err := w.registry.client.Get(w.ctx, prefix, clientv3.WithPrefix())
				if err != nil {
					continue
				}

				instances := make([]*registry.ServiceInstance, 0, len(resp.Kvs))
				for _, kv := range resp.Kvs {
					instance := &registry.ServiceInstance{}
					if err := json.Unmarshal(kv.Value, instance); err == nil {
						instances = append(instances, instance)
					}
				}

				select {
				case w.watchCh <- instances:
				case <-w.stopCh:
					return
				case <-w.ctx.Done():
					return
				}

				break
			}
		}
	}
}

// Next 监听服务实例变化
func (w *etcdWatcher) Next() ([]*registry.ServiceInstance, error) {
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
func (w *etcdWatcher) Stop() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.stopped {
		w.stopped = true
		w.cancel()
		close(w.stopCh)
	}

	return nil
}
