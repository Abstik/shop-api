package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mxshop-api/goods-web/models"
)

// 管理员权限校验
func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求上下文中获取set好的数据
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		// 如果不是2（管理员），则无权限
		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}

		// 如果是管理员，则继续
		ctx.Next()
	}
}
