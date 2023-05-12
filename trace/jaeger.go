package trace

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/whj1990/go-core/config"
)

func Init() io.Closer {
	return initGlobalTracer(
		config.GetNacosConfigData().Jaeger.ServiceName,
		config.GetNacosConfigData().Jaeger.HostPort,
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
	/*sender := transport.NewHTTPTransport(
		hostPort,
	)
	tracer, closer := jaeger.NewTracer(serviceName,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.Logger(jaeger.StdLogger)),
	)*/
	opentracing.InitGlobalTracer(tracer)
	return closer
}
