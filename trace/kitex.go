package trace

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	tracer "github.com/kitex-contrib/tracer-opentracing"
	"github.com/opentracing/opentracing-go"
)

var operationNameFunc = func(ctx context.Context) string {
	endpoint := rpcinfo.GetRPCInfo(ctx).To()
	return endpoint.ServiceName() + "::" + endpoint.Method()
}

func NewServerSuite() server.Suite {
	return tracer.NewServerSuite(opentracing.GlobalTracer(), operationNameFunc)
}

func NewClientSuite() client.Suite {
	return tracer.NewClientSuite(opentracing.GlobalTracer(), operationNameFunc)
}
