package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"mxshop-api/goods-web/global"
)

// 配置链路追踪组件，并为商品服务的所有路由配置此组件
func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 定义Jaeger配置
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst, // 设置了采样器的类型，此类型意味着所有的请求都会被采样并追踪
				Param: 1,                       // 表示所有请求都会被采样
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,                                                                                           // 将追踪的跨度记录到日志中
				LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port), // 指定Jaeger代理的地址
			},
			ServiceName: global.ServerConfig.JaegerInfo.Name,
		}

		// 使用Jaeger配置创建一个新的追踪器（tracer），并且同时获取一个关闭函数和一个错误信息（err）
		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}

		// 将指定的跟踪器设置为全局可用的跟踪器，从而在整个应用程序中启用分布式追踪功能
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()

		// 创建一个span，将http请求的url路径作为span的名称
		// 记录从这个请求开始到结束的所有操作，以便后续分析和监控
		startSpan := tracer.StartSpan(ctx.Request.URL.Path)
		defer startSpan.Finish()

		// 将tracer和parentSpan存储在gin.Context的gin上下文中，以便在处理请求时使用
		ctx.Set("tracer", tracer)
		ctx.Set("parentSpan", startSpan)
		ctx.Next()
	}
}
