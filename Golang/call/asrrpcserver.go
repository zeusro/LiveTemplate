// +build rpc

package main

import (
	"flag"
	"fmt"

	"call/internal/config"
	"call/internal/handler"
	"call/internal/svc"
	"call/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/asr.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// 注册 gRPC 服务
		pb.RegisterASRServiceServer(grpcServer, handler.NewASRHandler(ctx))

		if c.RestConf.Mode == "dev" {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting gRPC server at %s\n", c.RpcServerConf.ListenOn)
	logx.Infof("Starting gRPC server at %s", c.RpcServerConf.ListenOn)
	s.Start()
}
