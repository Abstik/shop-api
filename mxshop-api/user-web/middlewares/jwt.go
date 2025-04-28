package middlewares

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"mxshop-api/user-web/global"
	"mxshop-api/user-web/models"
)

// 验证jwt
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort()
			return
		}

		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			// 如果是token已过期错误
			if errors.Is(err, TokenExpired) {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "授权已过期",
				})
				c.Abort()
				return
			}
			// 如果是其他错误
			c.JSON(http.StatusUnauthorized, "未登陆")
			c.Abort()
			return
		}
		// 将关键信息放入上下文
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

// 封装jwt签名密钥的结构体
type JWT struct {
	SigningKey []byte
}

// 常见的错误
var (
	TokenExpired     = errors.New("token is expired")           // token已过期
	TokenNotValidYet = errors.New("token not active yet")       // token未生效
	TokenMalformed   = errors.New("that's not even a token")    // token格式错误
	TokenInvalid     = errors.New("couldn't handle this token") // 无法处理这个token
)

// 初始化jwt签名密钥的结构体
func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTInfo.SigningKey), // 可以设置过期时间
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析token
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	// 调用jwt.ParseWithClaims方法解析token
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		// 获取签名密钥
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			/*ve.Errors 是一个位掩码，记录了 token 解析过程中遇到的各种错误。
			jwt.ValidationErrorMalformed是一个常量，表示token格式不合法。
			二者做按位与（&）运算：
			如果结果不为0 (!= 0)，说明 ve.Errors 里包含了 ValidationErrorMalformed 这个错误。
			如果结果是0，说明没有格式错误。*/
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// 如果是token格式错误
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// 如果是过期，尝试刷新token
				newToken, refreshErr := j.RefreshToken(tokenString)
				if refreshErr != nil {
					return nil, TokenExpired
				}
				// 刷新成功后，再用新的token解析一遍
				return j.ParseToken(newToken)
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				// 如果是未生效
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	// 如果token不为nil
	if token != nil {
		// 将token的Claims断言为自定义的jwt结构体CustomClaims
		// 如果断言成功并且token的Valid字段为true，则返回claims
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		// 否则返回错误
		return nil, TokenInvalid
	} else {
		// 如果token为nil，则返回错误
		return nil, TokenInvalid
	}
}

// 更新token未被使用
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	// 强制设置 JWT 库内部用来判断时间的方法，把当前时间设定为 1970-01-01 00:00:00
	// 目的：让即使过期的 Token 也能通过验证
	// 因为正常解析时，JWT 会判断 exp（过期时间），如果 Token 已经过期，会直接抛错。但这里人为“冻结”了时间，绕过过期检测
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	// 调用jwt.ParseWithClaims方法解析token
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}

	// 将token的Claims断言为自定义的jwt结构体CustomClaims
	// 如果断言成功并且token的Valid字段为true，则返回claims
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		// 恢复JWT库内部用来判断时间的方法
		jwt.TimeFunc = time.Now
		// 更新claims中的过期时间(exp)字段，把有效期设置为1小时后
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		// 重新生成token
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
