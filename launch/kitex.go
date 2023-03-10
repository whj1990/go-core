package launch

import (
	"github.com/whj1990/go-core/config"
	"github.com/whj1990/go-core/trace"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"time"
)

func RpcServerOptions() []server.Option {
	return []server.Option{
		server.WithServiceAddr(utils.NewNetAddr("tcp", ":"+config.GetString("server.port", ""))),
		server.WithReadWriteTimeout(60 * time.Second),
		server.WithSuite(trace.NewServerSuite()),
	}
}

func RunServer(server server.Server) {
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
