package launch

import (
	"math"
	"net"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/whj1990/go-core/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func GrpcServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 4),
		grpc.MaxSendMsgSize(math.MaxInt32),
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
		grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer())),
	}
}
func RunGrpcServer(server *grpc.Server) {
	listener, err := net.Listen(config.GetString("server.network", "tcp"), config.GetString("server.address", ""))
	zap.L().Info("net.Listing", zap.String("port", config.GetString("server.address", "")))
	if err != nil {
		zap.L().Error(err.Error())
	}
	if err = server.Serve(listener); err != nil {
		zap.L().Error(err.Error())
	}

}
