package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	"mxshop-api/user-web/utils"
	"mxshop-api/user-web/utils/register/consul"
	myvalidator "mxshop-api/user-web/validator"
)

func main() {
	//1. 初始化logger（zap）
	initialize.InitLogger()

	//2. 初始化配置文件（viper）
	initialize.InitConfig()
	zap.S().Info(global.ServerConfig)

	//3. 初始化路由
	Router := initialize.Routers()

	//4. 初始化翻译器
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	//5. 初始化srv的连接
	initialize.InitSrvConn()

	// 自动加载环境变量，并根据环境类型（本地开发或线上环境）确定端口号
	// 如果是本地开发环境，端口号固定；如果是线上环境，则动态获取端口号
	viper.AutomaticEnv()
	// 如果是本地开发环境端口号固定
	debug := viper.GetBool("MXSHOP_DEBUG")
	if !debug { // 线上环境则获取空闲端口号并更新为服务器配置的端口号
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	// 注册自定义的验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok { // 检查当前的验证引擎是否是指定类型的validator.Validate
		// 注册自定义的校验手机号验证规则
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)

		// 注册自定义的错误提示
		_ = v.RegisterTranslation("mobile", global.Trans,
			func(ut ut.Translator) error {
				return ut.Add("mobile", "{0}非法的手机号码!", true)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("mobile", fe.Field())
				return t
			})
	}

	// consul服务注册中心
	// 向consul服务中心注册自己，以便其他服务发现和调用
	// 创建consul客户端连接
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	// 生成唯一的服务id
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	// 将当前此服务注册到consul
	err := registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}

	// 启动 HTTP 服务器，监听指定端口
	zap.S().Debugf("启动服务器, 端口： %d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功:")
	}
}
