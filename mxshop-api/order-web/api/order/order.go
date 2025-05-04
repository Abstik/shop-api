package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"

	"mxshop-api/order-web/api"
	"mxshop-api/order-web/forms"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/models"
	"mxshop-api/order-web/proto"
)

// 获取订单列表
func List(ctx *gin.Context) {
	// 获取jwt校验的参数
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")

	request := proto.OrderFilterRequest{}

	// 如果是管理员用户则返回所有的订单
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		// 如果是普通用户，则设置用户id
		request.UserId = int32(userId.(uint))
	}
	// 如果是管理员用户，id为零值，gorm进行查询时忽略零值直接全表查询

	// 获取分页数
	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	// 获取每页的数量
	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	// 查询订单列表
	rsp, err := global.OrderSrvClient.OrderList(context.WithValue(context.Background(), "ginContext", ctx), &request)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	orderList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		tmpMap := map[string]interface{}{}

		tmpMap["id"] = item.Id
		tmpMap["user"] = item.UserId
		tmpMap["order_sn"] = item.OrderSn
		tmpMap["pay_type"] = item.PayType
		tmpMap["status"] = item.Status
		tmpMap["post"] = item.Post
		tmpMap["total"] = item.Total
		tmpMap["address"] = item.Address
		tmpMap["name"] = item.Name
		tmpMap["mobile"] = item.Mobile
		tmpMap["add_time"] = item.AddTime

		orderList = append(orderList, tmpMap)
	}
	reMap["data"] = orderList
	ctx.JSON(http.StatusOK, reMap)
}

// 新建订单
func New(ctx *gin.Context) {
	// 获取参数
	orderForm := forms.CreateOrderForm{}
	if err := ctx.ShouldBindJSON(&orderForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}

	// 获取userID
	userId, _ := ctx.Get("userId")

	// 把当前的 gin.Context（变量名是ctx）塞入go语言的context.Context中，并指定键名为ginContext
	// 改造后的otgrpc源码中，通过context.Context获取gin.Context，通过gin.Context获取tracer和parentSpan
	rsp, err := global.OrderSrvClient.CreateOrder(context.WithValue(context.Background(), "ginContext", ctx), &proto.OrderRequest{
		UserId:  int32(userId.(uint)),
		Name:    orderForm.Name,
		Mobile:  orderForm.Mobile,
		Address: orderForm.Address,
		Post:    orderForm.Post,
	})
	if err != nil {
		zap.S().Errorw("新建订单失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 生成支付宝的支付url
	alipayUrl, _ := GenerateAlipayUrl(ctx, rsp.OrderSn, rsp.Total)

	ctx.JSON(http.StatusOK, gin.H{
		"id":         rsp.Id,
		"alipay_url": alipayUrl,
	})
}

// 查询订单详情
func Detail(ctx *gin.Context) {
	// 获取订单id参数
	id := ctx.Param("id")
	// 获取用户id
	userId, _ := ctx.Get("userId")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	// 封装订单id
	request := proto.OrderRequest{
		Id: int32(i),
	}

	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		// 如果非管理员，则封装userID
		request.UserId = int32(userId.(uint))
	}
	// 如果是管理员则request.UserId默认为零值，gorm查询时忽略零值

	rsp, err := global.OrderSrvClient.OrderDetail(context.WithValue(context.Background(), "ginContext", ctx), &request)
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{}
	reMap["id"] = rsp.OrderInfo.Id
	reMap["status"] = rsp.OrderInfo.Status
	reMap["user"] = rsp.OrderInfo.UserId
	reMap["post"] = rsp.OrderInfo.Post
	reMap["total"] = rsp.OrderInfo.Total
	reMap["address"] = rsp.OrderInfo.Address
	reMap["name"] = rsp.OrderInfo.Name
	reMap["mobile"] = rsp.OrderInfo.Mobile
	reMap["pay_type"] = rsp.OrderInfo.PayType
	reMap["order_sn"] = rsp.OrderInfo.OrderSn

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Goods {
		tmpMap := gin.H{
			"id":    item.GoodsId,
			"name":  item.GoodsName,
			"image": item.GoodsImage,
			"price": item.GoodsPrice,
			"nums":  item.Nums,
		}

		goodsList = append(goodsList, tmpMap)
	}
	reMap["goods"] = goodsList

	// 生成支付宝的支付url
	alipayUrl, _ := GenerateAlipayUrl(ctx, rsp.OrderInfo.OrderSn, rsp.OrderInfo.Total)

	reMap["alipay_url"] = alipayUrl
	ctx.JSON(http.StatusOK, reMap)
}

// 修改订单状态
func Update(ctx *gin.Context) {
	// 获取订单id参数
	orderIdStr := ctx.Param("orderIdStr")
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	orderForm := forms.UpdateOrderForm{}
	if err := ctx.ShouldBindJSON(&orderForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}

	// 修改订单状态
	_, err = global.OrderSrvClient.UpdateOrderStatus(context.WithValue(context.Background(), "ginContext", ctx), &proto.OrderStatus{
		Id:      int32(orderId),
		OrderSn: orderForm.OrderSn,
		Status:  orderForm.Status,
	})
	if err != nil {
		zap.S().Errorw("修改订单状态失败")
		api.HandleGrpcErrorToHttp(err, ctx)
	}
}

// 生成支付宝url函数
func GenerateAlipayUrl(ctx *gin.Context, orderSn string, total float32) (string, error) {
	// 生成支付宝的支付url
	client, err := alipay.New(global.ServerConfig.AliPayInfo.AppID, global.ServerConfig.AliPayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return "", err
	}
	err = client.LoadAliPayPublicKey(global.ServerConfig.AliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return "", err
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AliPayInfo.NotifyURL
	p.ReturnURL = global.ServerConfig.AliPayInfo.ReturnURL
	p.Subject = "慕学生鲜订单-" + orderSn
	p.OutTradeNo = orderSn
	p.TotalAmount = strconv.FormatFloat(float64(total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付url失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
	return url.String(), nil
}
