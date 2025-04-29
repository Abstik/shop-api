package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"mxshop-api/user-web/global"
)

// 使用viper读取配置文件，配置文件.yaml中存放nacos配置信息，通过配置连接到nacos服务器，并读取剩余配置信息

func GetEnvInfo(env string) bool {
	// 自动加载环境变量
	viper.AutomaticEnv()
	// 获取指定环境变量env的布尔值并返回
	return viper.GetBool(env)
}

func InitConfig() {
	// 1.通过 MXSHOP_DEBUG 环境变量决定加载 config-debug.yaml（开发环境）或 config-pro.yaml（生产环境）
	// 从环境变量 MXSHOP_DEBUG 中获取调试信息
	debug := GetEnvInfo("MXSHOP_DEBUG")
	// 设置配置文件前缀为 "config"
	configFilePrefix := "config"
	// 根据调试模式来决定读取不同的配置文件
	configFileName := fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix)
	}

	// 2.使用 Viper 加载本地配置文件，获取 nacos 所需的基础配置
	v := viper.New()
	// 文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// 这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("nacos配置信息: &+v", global.NacosConfig)

	// 3.初始化 Nacos 连接参数，从 nacos 中读取配置信息
	// 初始化nacos客户端
	sc := []constant.ServerConfig{ // Nacos服务器地址列表
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{ // 客户端行为配置（超时、日志目录等）
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{ // 创建与Nacos服务器的连接客户端
		"serverConfigs": sc, // Nacos服务器地址列表
		"clientConfig":  cc, // 客户端行为配置（超时、日志目录等）
	})
	/*根据 sc 中的服务器地址尝试建立连接。
	使用 cc 中的配置初始化客户端行为（如超时时间、缓存策略等）*/
	if err != nil {
		panic(err)
	}

	// 从 Nacos 服务器拉取指定配置内容
	/*客户端向 Nacos 服务器发送 HTTP 请求获取配置。
	Nacos 根据 DataId + Group 定位具体配置内容*/
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId, // 配置的唯一标识
		Group:  global.NacosConfig.Group}) // 配置分组
	if err != nil {
		panic(err)
	}

	// 将配置内容content解析到全局变量global.ServerConfig中
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	}
}
