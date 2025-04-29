package router

import (
	"github.com/gin-gonic/gin"

	"mxshop-api/goods-web/api/category"
	"mxshop-api/goods-web/middlewares"
)

// 分类相关路由组
func InitCategoryRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("categorys").Use(middlewares.Trace())
	{
		CategoryRouter.GET("", category.List)                                                            // 查询所有分类
		CategoryRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.Delete) // 删除分类
		CategoryRouter.GET("/:id", category.Detail)                                                      // 根据id查询当前分类及其子分类
		CategoryRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.New)          // 新建分类
		CategoryRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.Update)    // 修改分类
	}
}
