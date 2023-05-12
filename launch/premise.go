package launch

import (
	"io"

	"github.com/whj1990/go-core/config"
	"github.com/whj1990/go-core/trace"
	"go.uber.org/zap"
)

func InitPremise() (*zap.Logger, io.Closer) {
	config.NaCosInitConfigClient()
	logger := Init()
	closer := trace.Init()
	return logger, closer
}
