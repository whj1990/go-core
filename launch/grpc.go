package launch

import (
	"net"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/whj1990/go-core/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GrpcServerOptions() []grpc.ServerOption {
	cer, err := credentials.NewServerTLSFromFile("./cert/server.crt", "./cert/server.key")
	if err != nil {
		zap.L().Error(err.Error())
	}
	return []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 5),
		grpc.MaxSendMsgSize(1024 * 1024 * 5),
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
		grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer())),
		grpc.Creds(cer),
	}
}
func RunGrpcServer(server *grpc.Server) {
	listener, err := net.Listen(config.GetNacosConfigData().GrpcServer.Network, config.GetNacosConfigData().GrpcServer.Address)
	zap.L().Info("net.Listing", zap.String("port", config.GetNacosConfigData().GrpcServer.Address))
	if err != nil {
		zap.L().Error(err.Error())
	}
	if err = server.Serve(listener); err != nil {
		zap.L().Error(err.Error())
	}

}
