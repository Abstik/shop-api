package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mxshop-api/oss-web/utils"
)

func Token(c *gin.Context) {
	response := utils.GetPolicyToken()
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.String(200, response)
}

func Callback(ctx *gin.Context) {
	// 读取 OSS 回调发送的表单数据
	filename := ctx.PostForm("filename")
	size := ctx.PostForm("size")
	mimeType := ctx.PostForm("mimeType")
	height := ctx.PostForm("height")
	width := ctx.PostForm("width")

	zap.S().Infof("OSS 回调接收成功：filename=%s, size=%s, mimeType=%s, height=%s, width=%s",
		filename, size, mimeType, height, width)

	// 你可以在这里做你自己的业务逻辑，比如保存到数据库，发送消息等等

	// 响应给阿里云，必须是 JSON 格式，否则阿里云认为回调失败
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
	})
}
