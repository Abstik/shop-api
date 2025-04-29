package forms

// 密码登录表单参数
type PassWordLoginForm struct {
	Email     string `form:"email" json:"email" binding:"required,email"`
	PassWord  string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"` // 图片验证码
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`       // 图片验证码id
}

// 注册表单参数
type RegisterForm struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Code     string `form:"code" json:"code" binding:"required,min=6,max=6"` // 手机短信验证码
}

// 修改用户信息表单参数
type UpdateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=3,max=10"`
	Gender   string `form:"gender" json:"gender" binding:"required,oneof=female male"`
	Birthday string `form:"birthday" json:"birthday" binding:"required,datetime=2006-01-02"`
}
