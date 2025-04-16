package models

import (
	"github.com/dgrijalva/jwt-go"
)

// 自定义jwt结构体
type CustomClaims struct {
	ID                 uint
	NickName           string // 昵称
	AuthorityId        uint   // 角色id
	jwt.StandardClaims        // 内嵌jwt标准结构体
}
