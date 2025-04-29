package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/router"
)

// 初始化路由

func Routers() *gin.Engine {
	Router := gin.New()
	// 添加Recovery中间件
	Router.Use(gin.Recovery())
	// 添加自定义Logger中间件，跳过/health路径
	Router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))

	// 配置健康检查
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	// 配置跨域
	Router.Use(middlewares.Cors())

	// 创建总路由组
	ApiGroup := Router.Group("/u/v1")

	// 初始化用户路由组
	router.InitUserRouter(ApiGroup)
	// 初始化基础路由组(验证码相关)
	router.InitBaseRouter(ApiGroup)

	return Router
}
