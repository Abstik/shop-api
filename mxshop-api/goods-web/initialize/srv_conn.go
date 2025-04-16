package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver" // 导包但不使用
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
)

// 初始化grpc连接
func InitSrvConn() {
	// 运用从consul获取的服务配置
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		//grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	// 创建grpc客户端
	global.GoodsSrvClient = proto.NewGoodsClient(userConn)
}
