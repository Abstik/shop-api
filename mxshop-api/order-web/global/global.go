package global

import (
	ut "github.com/go-playground/universal-translator"

	"mxshop-api/order-web/config"
	"mxshop-api/order-web/proto"
)

// 全局变量
var (
	Trans ut.Translator

	ServerConfig = &config.ServerConfig{}

	NacosConfig = &config.NacosConfig{}

	// 商品服务
	GoodsSrvClient proto.GoodsClient

	// 订单服务
	OrderSrvClient proto.OrderClient

	// 库存服务
	InventorySrvClient proto.InventoryClient
)
