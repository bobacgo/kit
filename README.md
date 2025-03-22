# kit

![](./image.png)

### A microservice framework

example  [bobacgo/go-wab](https://github.com/bobacgo/go-wab)

### Dependency Package

* gin
* grpc
* zap
* gopkg.in/yaml.v3
* validator
* redis
* gorm
* mysql

### Install

```shell
go get github.com/bobacgo/kit
```

### Quick Start

```go
package main

import (
	"flag"
	kserver "github.com/bobacgo/kit/app/server"
	"github.com/bobacgo/kit/examples/internal/server"
	"log"

	"github.com/bobacgo/kit/app"
	"github.com/bobacgo/kit/app/conf"
	"github.com/bobacgo/kit/examples/config"
	"github.com/bobacgo/kit/examples/internal/app/router"
)

var filepath = flag.String("config", "./config.yaml", "config file path")

func init() {
	flag.String("name", "admin-service", "service name")
	flag.String("env", "dev", "run config context")
	flag.String("logger.level", "info", "logger level")
	flag.Int("port", 8080, "http port 8080, rpc port 9080")
	conf.BindPFlags()
}

func main() {
	newApp := app.New[config.Service](*filepath,
		// app.WithMustDB(),
		// app.WithMustRedis(),
		app.WithGinServer(router.Register),
		app.WithServer("kafka", func(a *app.Options) kserver.Server {
			return new(server.KafkaServer)
		}),
	)
	if err := newApp.Run(); err != nil {
		log.Panic(err.Error())
	}
}

```

### app run log

```plaintext
go run main.go --config ./config.yaml
2025-03-06 14:00:22.870 INFO    app/option.go:116       local config info
basic:
  name: examples-service
  version: 1.0.0
  env: dev
  configs:
  - ./deploy/v1.0.0/db.yaml
  - ./deploy/v1.0.0/logger.yaml
  - ./deploy/v1.0.0/redis.yaml
  registry:
    addr: 127.0.0.1:2379
    timeout: ""
  server:
    http:
      addr: 0.0.0.0:8080
      timeout: 1s
    rpc:
      addr: 0.0.0.0:9080
      timeout: 1s
  security:
    ciphertext:
      isCiphertext: false
      cipherKey: YpC5w*****uMvd4f
    jwt:
      secret: YpC5w*****uMvd4f
      issuer: gogo-admin
      accessTokenExpired: ""
      refreshTokenExpired: ""
      cacheKeyPrefix: admin:login_token
  logger:
    level: ""
    timeFormat: "2006-01-02 15:04:05.000"
    filepath: ./logs
    filename: examples-service
    filenameSuffix: 2006-01-02-150405
    fileExtension: log
    fileJsonEncoder: false
    fileSizeMax: 10
    fileAgeMax: 180
    fileCompress: true
  db:
    default:
      source: admin******tcp(127.0.0.1:3306)/mall-ums?charset=utf8mb4&parseTime=True&loc=Local
      dryRun: false
      slowThreshold: 1
      maxLifeTime: 1
      maxOpenConn: 100
      maxIdleConn: 30
  localCache:
    maxSize: 500MB
  redis:
    addrs:
    - 127.0.0.1:6379
    username: ""
    password: ""
    db: 0
    poolSize: 50
    readTimeout: 1ms
    writeTimeout: 1ms
service:
  errattemptlimit: 5

2025-03-06 14:00:22.870 INFO    app/option.go:118       [config] init done.
2025-03-06 14:00:22.870 INFO    app/option.go:119       [logger] init done.
2025-03-06 14:00:22.877 INFO    app/option.go:132       [local_cache] init done.
2025-03-06 14:00:22.877 INFO    app/service.go:95       server started
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

2025-03-06 14:00:22.878 WARN    app/http.go:45  [gin] Running in "debug" mode
[GIN-debug] GET    /health                   --> github.com/bobacgo/kit/app.healthApi.func1 (4 handlers)
[GIN-debug] GET    /debug/pprof/             --> github.com/bobacgo/kit/app.pprofApi.WrapF.func1 (4 handlers)
[GIN-debug] GET    /debug/pprof/cmdline      --> github.com/bobacgo/kit/app.pprofApi.WrapF.func2 (4 handlers)
[GIN-debug] GET    /debug/pprof/profile      --> github.com/bobacgo/kit/app.pprofApi.WrapF.func3 (4 handlers)
[GIN-debug] GET    /debug/pprof/symbol       --> github.com/bobacgo/kit/app.pprofApi.WrapF.func4 (4 handlers)
[GIN-debug] GET    /debug/pprof/trace        --> github.com/bobacgo/kit/app.pprofApi.WrapF.func5 (4 handlers)
[GIN-debug] GET    /debug/pprof/allocs       --> github.com/bobacgo/kit/app.pprofApi.WrapF.func6 (4 handlers)
[GIN-debug] GET    /debug/pprof/block        --> github.com/bobacgo/kit/app.pprofApi.WrapF.func7 (4 handlers)
[GIN-debug] GET    /debug/pprof/goroutine    --> github.com/bobacgo/kit/app.pprofApi.WrapF.func8 (4 handlers)
[GIN-debug] GET    /debug/pprof/heap         --> github.com/bobacgo/kit/app.pprofApi.WrapF.func9 (4 handlers)
[GIN-debug] GET    /debug/pprof/mutex        --> github.com/bobacgo/kit/app.pprofApi.WrapF.func10 (4 handlers)
[GIN-debug] GET    /debug/pprof/threadcreate --> github.com/bobacgo/kit/app.pprofApi.WrapF.func11 (4 handlers)
[GIN-debug] POST   /v1/user/create           --> github.com/bobacgo/kit/examples/internal/app/admin.Register.func1 (4 handlers)
[GIN-debug] PUT    /v1/user/update           --> github.com/bobacgo/kit/examples/internal/app/admin.Register.func2 (4 handlers)
[GIN-debug] DELETE /v1/user/delete           --> github.com/bobacgo/kit/examples/internal/app/admin.Register.func3 (4 handlers)
[GIN-debug] POST   /v1/user/pageList         --> github.com/bobacgo/kit/examples/internal/app/admin.Register.func4 (4 handlers)
2025-03-06 14:00:22.878 INFO    app/http.go:61  http server running http://192.168.1.4:8080
```