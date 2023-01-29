package launch

import (
	"github.com/whj1990/go-core/config"
	"github.com/whj1990/go-core/trace"
	"go.uber.org/zap"
	"io"
)

func InitPremise() (*zap.Logger, io.Closer) {
	config.Init()
	logger := Init()
	closer := trace.Init()
	return logger, closer
}
