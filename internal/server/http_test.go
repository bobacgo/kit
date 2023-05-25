package server_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/gogo/internal/server"
)

func Test_HttpServer(t *testing.T) {
	exitHttp := make(chan struct{})
	doneExitHttp := make(chan struct{})
	server.RunHttpServer(exitHttp, doneExitHttp, ":8080", router)
}

func Test_HttpApi(t *testing.T) {
	r, err := http.DefaultClient.Get("http://127.0.0.1:8080/ping")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		t.Fatalf("http status: %s", r.Status)
	}
	b, err2 := io.ReadAll(r.Body)
	if err2 != nil {
		t.Fatal(err2)
	}
	t.Logf("%s", b)
}

func router(e *gin.Engine) {
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			"code": 0,
			"msg":  "ok",
		})
	})
}
