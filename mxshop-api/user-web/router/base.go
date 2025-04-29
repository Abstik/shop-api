package router

import (
	"github.com/gin-gonic/gin"

	"mxshop-api/user-web/api"
)

// 验证码相关路由
func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		// 登录时生成图片验证码
		BaseRouter.GET("captcha", api.GetCaptcha)
		// 注册时发送短信验证码
		BaseRouter.POST("send_sms", api.SendVerificationCode)
	}

}
