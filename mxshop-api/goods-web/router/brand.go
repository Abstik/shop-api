package router

import (
	"github.com/gin-gonic/gin"

	"mxshop-api/goods-web/api/brands"
	"mxshop-api/goods-web/middlewares"
)

func InitBrandRouter(Router *gin.RouterGroup) {
	// 品牌相关路由组
	BrandRouter := Router.Group("brands").Use(middlewares.Trace())
	{
		BrandRouter.GET("", brands.BrandList)                                                            // 品牌列表页
		BrandRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), brands.DeleteBrand) // 删除品牌
		BrandRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), brands.NewBrand)          // 新建品牌
		BrandRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), brands.UpdateBrand)    // 修改品牌信息
		//BrandRouter.GET("/:id", brands.GetBrand)                                                         // 根据id查询品牌信息
	}

	// 品牌和分类相关路由组
	CategoryBrandRouter := Router.Group("categorybrands")
	{
		CategoryBrandRouter.GET("", brands.CategoryBrandList)                                                            // 查询所有品牌和类别的关联信息
		CategoryBrandRouter.GET("/:id", brands.GetCategoryBrandList)                                                     // 根据类别查询此类别下的所有品牌
		CategoryBrandRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), brands.NewCategoryBrand)          // 新建品牌和类别的关联
		CategoryBrandRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), brands.DeleteCategoryBrand) // 删除品牌和类别的关联
		CategoryBrandRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), brands.UpdateCategoryBrand)    // 更新品牌和类别的关联
	}
}
