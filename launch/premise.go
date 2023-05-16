package launch

import (
	"io"

	"github.com/whj1990/go-core/config"
	"github.com/whj1990/go-core/trace"
	"go.uber.org/zap"
)

func InitPremise(isServer bool) (*zap.Logger, io.Closer) {
	config.NaCosInitConfigClient()
	config.NewNaCosNamingClient(isServer)
	logger := Init()
	closer := trace.Init()
	return logger, closer
}
