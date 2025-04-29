package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mxshop-api/order-web/middlewares"
	"mxshop-api/order-web/router"
)

func Routers() *gin.Engine {
	Router := gin.New()
	// 添加Recovery中间件
	Router.Use(gin.Recovery())
	// 添加自定义Logger中间件，跳过/health路径
	Router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))

	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	//配置跨域
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/o/v1")

	// 初始化订单路由组
	router.InitOrderRouter(ApiGroup)
	// 初始化购物车路由组
	router.InitShopCartRouter(ApiGroup)

	return Router
}
