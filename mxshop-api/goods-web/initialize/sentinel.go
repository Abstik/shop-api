package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

// 初始化InitSentinel，进行熔断限流
func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalf("初始化sentinel 异常: %v", err)
	}

	// 配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "goods-list",
			TokenCalculateStrategy: flow.Direct, // 直接使用规则中的 Threshold 表示当前统计周期内的最大Token数量
			ControlBehavior:        flow.Reject, // 请求数超过了阈值，就直接拒绝
			Threshold:              20,          // 最多访问次数
			StatIntervalInMs:       6000,        // 统计周期，单位毫秒
		},
	})

	if err != nil {
		zap.S().Fatalf("加载规则失败: %v", err)
	}
}
