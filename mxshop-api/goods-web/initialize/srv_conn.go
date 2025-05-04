package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
	"mxshop-api/goods-web/utils/otgrpc"
)

// 初始化grpc连接
func InitSrvConn() {
	// 运用从consul获取的服务配置
	consulInfo := global.ServerConfig.ConsulInfo

	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),                    // 指定负载均衡策略为“轮询”（round robin）
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())), // 添加OpenTracing客户端拦截器
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 商品服务失败")
	}
	// 创建grpc客户端
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	inventoryClient, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),                    // 指定负载均衡策略为“轮询”（round robin）
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())), // 添加OpenTracing客户端拦截器
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 库存服务失败", err)
	}

	// 创建grpc客户端
	global.InventorySrvClient = proto.NewInventoryClient(inventoryClient)
}
