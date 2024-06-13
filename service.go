package kit

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bobacgo/kit/pkg/network"
	"github.com/bobacgo/kit/pkg/uid"
	"golang.org/x/sync/errgroup"
)

type App struct {
	opts Options

	wg     sync.WaitGroup
	exit   chan struct{}
	signal chan os.Signal // 终止服务的信号
}

func New(opts ...Option) *App {
	wg, _ := errgroup.WithContext(context.Background())
	o := Options{
		appid: uid.UUID(),
		sigs:  []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		wg:    wg,
	}
	for _, opt := range opts {
		opt(&o)
	}
	if err := o.wg.Wait(); err != nil { // 等待 options 实例化结束
		log.Panic(err)
	}
	return &App{
		opts:   o,
		exit:   make(chan struct{}),
		signal: make(chan os.Signal, 1),
	}
}

func (a *App) Run() error {
	opts := a.opts
	if opts.httpServer != nil {
		go RunMustHttpServer(a, opts.httpServer)
	}
	if opts.rpcServer != nil {
		go RunMustRpcServer(a, opts.rpcServer)
	}

	ctx := context.Background()

	slog.Info("server started")
	// 阻塞并监听退出信号
	signal.Notify(a.signal, opts.sigs...)
	<-a.signal

	err := a.shutdown(ctx)
	slog.Info("service has exited")
	return err
}

// Stop 手动停止服务
func (a *App) Stop() {
	a.signal <- os.Interrupt
}

func (a *App) shutdown(ctx context.Context) error {
	// opts := a.opts
	close(a.exit) // 通知http、rpc服务退出信号

	// 1.等待 Http 服务结束退出
	// 2.等待 RPC 服务结束退出
	a.wg.Wait()
	return nil
}

func getRegistryUrl(scheme, addr string) (string, error) {
	ip, err := network.OutBoundIP()
	if err != nil {
		return "", err
	}
	_, ports, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	return scheme + "://" + net.JoinHostPort(ip, ports), nil
}
