package trace

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

func NewGormLogger() gormLogger.Interface {
	return &logger{}
}

type logger struct {
	gormLogger.Writer
	gormLogger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func (l *logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {}

func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {}

func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {}

func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	span, newCtx := opentracing.StartSpanFromContext(
		ctx,
		fmt.Sprintf("mysql::%s", sql),
		opentracing.Tag{Key: "rows", Value: rows},
		opentracing.Tag{Key: "err", Value: err},
	)
	opentracing.SpanFromContext(newCtx)
	span.Finish()
}
