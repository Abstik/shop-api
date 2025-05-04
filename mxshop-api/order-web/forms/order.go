package forms

// 新建订单表单参数
type CreateOrderForm struct {
	Address string `json:"address" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Mobile  string `json:"mobile" binding:"required"`
	Post    string `json:"post" binding:"required"`
}

// 修改订单状态表单参数
type UpdateOrderForm struct {
	OrderSn string `json:"order_sn" binding:"required"`
	Status  string `json:"status" binding:"required"`
}
