package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mxshop-api/userop-web/middlewares"
	"mxshop-api/userop-web/router"
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

	ApiGroup := Router.Group("/up/v1")
	router.InitAddressRouter(ApiGroup)
	router.InitMessageRouter(ApiGroup)
	router.InitUserFavRouter(ApiGroup)

	return Router
}
