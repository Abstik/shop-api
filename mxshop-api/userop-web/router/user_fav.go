package router

import (
	"github.com/gin-gonic/gin"

	"mxshop-api/userop-web/api/user_fav"
	"mxshop-api/userop-web/middlewares"
)

func InitUserFavRouter(Router *gin.RouterGroup) {
	UserFavRouter := Router.Group("userfavs").Use(middlewares.JWTAuth())
	{
		UserFavRouter.GET("", user_fav.List)          // 查询收藏信息
		UserFavRouter.POST("", user_fav.New)          // 新建收藏记录
		UserFavRouter.DELETE("/:id", user_fav.Delete) // 删除收藏记录
		UserFavRouter.GET("/:id", user_fav.Detail)    // 查看用户是否收藏过某件商品
	}
}
