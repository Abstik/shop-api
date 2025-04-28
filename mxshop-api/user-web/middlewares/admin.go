package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mxshop-api/user-web/models"
)

// 检验管理员权限
func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		// 只有管理员(角色id为2)才有权限
		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}

}
