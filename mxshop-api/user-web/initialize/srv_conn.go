package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

// 初始化grpc连接池
func InitSrvConn() {
	// 获取consul配置文件实例
	consulInfo := global.ServerConfig.ConsulInfo

	// 通过grpc连接到consul服务中心，通过服务发现机制访问已经注册的user_srv服务
	/*通过global.ServerConfig.UserSrvInfo.Name这个服务名，consul会根据这个服务名查找相关的服务实例，然后返回其可用的地址和端口 ，客户端就可以连接到该服务*/
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 提供负载均衡策略的配置。这里使用的是 round_robin，即轮询负载均衡，确保请求均匀分配到多个服务实例
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	// 创建连接（用户服务客户端）
	userSrvClient := proto.NewUserClient(userConn)
	// 赋值给全局变量
	global.UserSrvClient = userSrvClient
}
