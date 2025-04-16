package router

import (
	"github.com/gin-gonic/gin"

	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

//用户相关路由

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	{
		// 获取用户列表（只有管理员才有权限）
		UserRouter.GET("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		// 用户登录
		UserRouter.POST("pwd_login", api.PassWordLogin)
		// 用户注册
		UserRouter.POST("register", api.Register)
		// 查询用户信息
		UserRouter.GET("detail", middlewares.JWTAuth(), api.GetUserDetail)
		// 修改用户信息
		UserRouter.PATCH("update", middlewares.JWTAuth(), api.UpdateUser)
	}
	//服务注册和发现
}
