package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"

	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
)

// 生成6位随机验证码
func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// 发送验证码到邮箱并存入缓存
func SendVerificationCode(ctx *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// 配置发件人邮箱
	smtpHost := "smtp.qq.com"
	smtpPort := 465
	fromEmail := "2455494167@qq.com"
	authCode := "rpqcsjeyqoesecbd"

	// 生成验证码
	verificationCode := generateVerificationCode()

	// 创建邮件
	m := gomail.NewMessage()
	// m.FormatAddress将发件人邮箱和名称格式化为MIME编码的地址
	m.SetHeader("From", m.FormatAddress(fromEmail, "农视界"))
	m.SetHeader("To", sendSmsForm.Email)
	m.SetHeader("Subject", "验证码")
	m.SetBody("text/html", fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>验证码</title>
			<style>
				body { font-family: Arial, sans-serif; }
				.container { padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
				h1 { color: #333; }
				.code { font-size: 24px; font-weight: bold; color: #007bff; }
				.footer { margin-top: 20px; font-size: 12px; color: #888; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>你的验证码</h1>
				<p class="code">%s</p>
				<p>有效时间为 5 分钟</p>
				<div class="footer">如果您没有请求此验证码，请忽略此邮件。</div>
			</div>
		</body>
		</html>
	`, verificationCode))

	// 发送邮件
	d := gomail.NewDialer(smtpHost, smtpPort, fromEmail, authCode)
	d.SSL = true
	if err := d.DialAndSend(m); err != nil {
		zap.S().Errorf("发送邮件失败: %v", err)
	}

	// 将验证码保存起来 - redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	rdb.Set(context.Background(), sendSmsForm.Email, verificationCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
