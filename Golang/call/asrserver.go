// +build !rpc

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"call/internal/config"
	"call/internal/handler"
	"call/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/asr.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)

	// 创建 HTTP 服务器
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册 WebSocket 路由
	wsHandler := handler.NewWebSocketHandler(ctx)
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/ws",
		Handler: wsHandler.HandleWebSocket,
	})

	// 注册健康检查
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/health",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	})

	// 注册会话查询接口
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/session/:sessionId",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			sessionID := r.PathValue("sessionId")
			session, err := ctx.SessionMgr.GetSession(r.Context(), sessionID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(session)
		},
	})

	fmt.Printf("Starting HTTP server at %s:%d\n", c.Host, c.Port)
	logx.Infof("Starting HTTP server at %s:%d", c.Host, c.Port)
	server.Start()
}
