package global

import (
	ut "github.com/go-playground/universal-translator"

	"mxshop-api/user-web/config"
	"mxshop-api/user-web/proto"
)

// 全局变量

var (
	// 翻译器
	Trans ut.Translator

	// 全局配置
	ServerConfig = &config.ServerConfig{}

	// nacos配置
	NacosConfig = &config.NacosConfig{}

	// 用户服务的客户端，可以调用相应的注册的用户服务
	UserSrvClient proto.UserClient
)
