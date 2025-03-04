# gogo

### A fully functional web framework

example  [gogoclouds/go-wab](https://github.com/gogoclouds/go-wab)

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
	newApp := app.New(
		app.WithMustConfig(*filepath, func(cfg *conf.ServiceConfig[config.Service]) {
			config.Cfg = &cfg.Service
		}),
		app.WithLogger(),
		app.WithMustLocalCache(),
		app.WithMustDB(),
		app.WithMustRedis(),
		app.WithGinServer(router.Register),
	)
	if err := newApp.Run(); err != nil {
		log.Panic(err.Error())
	}
}

```

### app run log

```plaintext
lanjin@lovec examples % go run main.go
2024/06/17 23:01:32 INFO local config info
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
      cipherKey: YpC5wIRf4ZuMvd4f
    jwt:
      secret: YpC5wIRf4ZuMvd4f
      issuer: gogo-admin
      accessTokenExpired: ""
      refreshTokenExpired: ""
      cacheKeyPrefix: admin:login_token
  logger:
    level: debug
    timeFormat: "2006-01-02 15:04:05.000"
    filepath: ./logs
    filename: ""
    filenameSuffix: 2006-01-02-150405
    fileExtension: log
    fileJsonEncoder: true
    fileSizeMax: 10
    fileAgeMax: 30
    fileCompress: true
  db:
    default:
      source: root:root@tcp(127.0.0.1:3306)/mall-ums?charset=utf8mb4&parseTime=True&loc=Local
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
  errAttemptLimit: 5
2024-06-17 23:01:32.833 INFO    app/option.go:121       [config] init done.
2024-06-17 23:01:32.833 INFO    app/option.go:122       [logger] init done.
2024-06-17 23:01:32.837 INFO    app/option.go:150       [redis] init done.
2024-06-17 23:01:32.839 INFO    app/option.go:172       [database] init done.
2024-06-17 23:01:32.842 INFO    app/option.go:135       [local_cache] init done.
2024-06-17 23:01:32.842 INFO    app/service.go:83       server started
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

2024-06-17 23:01:32.842 WARN    app/http.go:44  [gin] Running in "debug" mode
[GIN-debug] GET    /health                   --> github.com/bobacgo/kit/app.healthApi.func1 (4 handlers)
[GIN-debug] POST   /v1/user/create           --> github.com/bobacgo/kit/examples/internal/app/admin.Register.func1 (4 handlers)
[GIN-debug] PUT    /v1/user/update           --> github.com/bobacgo/kit/examples/internal/app/admin.Register.func2 (4 handlers)
[GIN-debug] DELETE /v1/user/delete           --> github.com/bobacgo/kit/examples/internal/app/admin.Register.func3 (4 handlers)
[GIN-debug] POST   /v1/user/pageList         --> github.com/bobacgo/kit/examples/internal/app/admin.Register.func4 (4 handlers)
2024-06-17 23:01:32.843 INFO    app/http.go:59  http server running http://192.168.1.3:8080
[GIN] 2024/06/17 - 23:02:12 | 404 |      21.375µs |             ::1 | GET      "/web/"
[GIN] 2024/06/17 - 23:02:12 | 404 |       1.167µs |             ::1 | GET      "/favicon.ico"
^C2024-06-17 23:02:43.171       INFO    app/http.go:67  Shutting down http server...
2024-06-17 23:02:43.172 INFO    app/http.go:75  http server exiting
2024-06-17 23:02:43.172 INFO    app/service.go:89       service has exited
```
