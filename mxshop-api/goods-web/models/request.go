package models

import (
	"github.com/dgrijalva/jwt-go"
)

// 自定义Claims结构体（jwt）
type CustomClaims struct {
	ID          uint
	NickName    string // 昵称
	AuthorityId uint   // 权限id
	jwt.StandardClaims
}
