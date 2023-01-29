package middleware

import (
	"fmt"
	"github.com/whj1990/go-core/constant"
	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"golang.org/x/net/context"
)

func ContextMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var span opentracing.Span
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
		if err != nil {
			span = opentracing.StartSpan(
				fmt.Sprintf("%s:%s", ctx.Request.Method, ctx.Request.URL.Path),
				ext.SpanKindRPCServer,
			)
		} else {
			span = opentracing.StartSpan(
				fmt.Sprintf("%s:%s", ctx.Request.Method, ctx.Request.URL.Path),
				opentracing.ChildOf(spCtx),
				ext.SpanKindRPCServer,
			)
		}
		defer span.Finish()
		opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
		ctx.Set(constant.RPCContext, setMetaInfo(ctx, span))
		ctx.Next()
	}
}

func setMetaInfo(ctx *gin.Context, span opentracing.Span) context.Context {
	rpcContext := opentracing.ContextWithSpan(ctx, span)
	rpcContext = metainfo.WithValue(rpcContext, "token", ctx.Request.Header.Get("Authorization"))
	rpcContext = metainfo.WithPersistentValue(rpcContext, constant.CurrentUserId, ctx.Request.Header.Get(constant.CurrentUserId))
	rpcContext = metainfo.WithPersistentValue(rpcContext, constant.CurrentUserName, ctx.Request.Header.Get(constant.CurrentUserName))
	rpcContext = metainfo.WithPersistentValue(rpcContext, constant.CurrentOrganizationId, ctx.Request.Header.Get(constant.CurrentOrganizationId))
	rpcContext = metainfo.WithPersistentValue(rpcContext, constant.AuthElements, ctx.Request.Header.Get(constant.AuthElements))
	return rpcContext
}
