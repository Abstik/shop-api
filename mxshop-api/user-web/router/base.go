package router

import (
	"github.com/gin-gonic/gin"

	"mxshop-api/user-web/api"
)

//验证码相关路由

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		// 生成验证码
		BaseRouter.GET("captcha", api.GetCaptcha)
		// 发送短信验证码
		BaseRouter.POST("send_sms", api.SendSms)
	}

}
