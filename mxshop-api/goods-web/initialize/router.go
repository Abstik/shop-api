package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mxshop-api/goods-web/middlewares"
	"mxshop-api/goods-web/router"
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

	//添加链路追踪
	ApiGroup := Router.Group("/g/v1")
	router.InitBannerRouter(ApiGroup)
	router.InitBrandRouter(ApiGroup)
	router.InitCategoryRouter(ApiGroup)
	router.InitGoodsRouter(ApiGroup)

	return Router
}
