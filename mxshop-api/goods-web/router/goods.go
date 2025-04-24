package router

import (
	"github.com/gin-gonic/gin"

	"mxshop-api/goods-web/api/goods"
	"mxshop-api/goods-web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	// 创建商品路由组并配置链路追踪
	// 每次调用商品路由组下的路由时，都会经过链路追踪中间件
	GoodsRouter := Router.Group("goods").Use(middlewares.Trace())
	{
		// 条件查询商品列表（url参数）
		GoodsRouter.GET("", goods.List)
		// 新增商品（需要管理员权限）
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)
		// 查询商品详情
		GoodsRouter.GET("/:id", goods.Detail)
		// 删除商品（需要管理员权限）
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)
		// 获取商品的库存
		GoodsRouter.GET("/:id/stocks", goods.Stocks)
		// 修改商品信息（需要管理员权限）
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
		// 修改商品状态（需要管理员权限）
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)
	}
}
