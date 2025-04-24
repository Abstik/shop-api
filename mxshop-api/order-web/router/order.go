package router

import (
	"github.com/gin-gonic/gin"

	"mxshop-api/order-web/api/order"
	"mxshop-api/order-web/api/pay"
	"mxshop-api/order-web/middlewares"
)

// 初始化订单路由
func InitOrderRouter(Router *gin.RouterGroup) {
	// 订单相关路由组，并配置链路追踪
	OrderRouter := Router.Group("orders").Use(middlewares.JWTAuth()).Use(middlewares.Trace())
	{
		OrderRouter.GET("", order.List)       // 订单列表
		OrderRouter.POST("", order.New)       // 新建订单
		OrderRouter.GET("/:id", order.Detail) // 订单详情
	}

	// 支付宝回调通知（当用户在支付宝完成支付后，支付宝会调用这个接口给服务端进行通知）
	PayRouter := Router.Group("pay")
	{
		PayRouter.POST("alipay/notify", pay.Notify)
	}
}
