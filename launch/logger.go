package launch

import (
	"github.com/whj1990/go-core/os"
	"go.uber.org/zap"
)

func Init() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	if os.RunningInDocker() {
		logger, _ = zap.NewProduction()
	}
	zap.ReplaceGlobals(logger)
	return logger
}
