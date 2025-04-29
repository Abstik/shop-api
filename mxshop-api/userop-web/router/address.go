package router

import (
	"github.com/gin-gonic/gin"

	"mxshop-api/userop-web/api/address"
	"mxshop-api/userop-web/middlewares"
)

func InitAddressRouter(Router *gin.RouterGroup) {
	AddressRouter := Router.Group("address").Use(middlewares.JWTAuth())
	{
		AddressRouter.GET("", address.List)
		AddressRouter.POST("", address.New)
		AddressRouter.DELETE("/:id", address.Delete)
		AddressRouter.PUT("/:id", address.Update)
	}
}
