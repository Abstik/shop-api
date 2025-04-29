package forms

type SendSmsForm struct {
	Email string `form:"email" json:"email" binding:"required,email"`
	Type  uint   `form:"type" json:"type" binding:"required,oneof=1 2"` // 1 表示注册 2表示找回密码
}
