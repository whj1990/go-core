package trace

import (
	"fmt"
	"github.com/whj1990/go-core/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
)

func Init() io.Closer {
	return initGlobalTracer(
		config.GetString("jaeger.serviceName", ""),
		config.GetString("jaeger.hostPort", ""),
	)
}

func initGlobalTracer(serviceName, hostPort string) io.Closer {
	cfg := &jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const", //固定采样
			Param: 1,       //1=全采样、0=不采样
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           false,
			LocalAgentHostPort: hostPort,
		},
		ServiceName: serviceName,
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.InitGlobalTracer(tracer)
	return closer
}
