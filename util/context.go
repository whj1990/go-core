package util

import (
	"encoding/json"
	"github.com/whj1990/go-core/constant"
	"github.com/whj1990/go-core/handler"
	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"strconv"
)

type OrganizationData struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func GetRPCContext(ctx *gin.Context) context.Context {
	rpcContext, _ := ctx.Get(constant.RPCContext)
	return rpcContext.(context.Context)
}

func GetCurrentUserId(ctx *gin.Context) (int64, error) {
	currentUserId, err := strconv.ParseInt(ctx.Request.Header.Get(constant.CurrentUserId), 10, 0)
	if err != nil {
		return 0, handler.HandleError(err)
	}
	return currentUserId, nil
}

func GetCurrentOrganizationId(ctx *gin.Context) (int64, error) {
	currentOrganizationId, err := strconv.ParseInt(ctx.Request.Header.Get(constant.CurrentOrganizationId), 10, 0)
	if err != nil {
		return 0, handler.HandleError(err)
	}
	return currentOrganizationId, nil
}

func GetMetaInfoToken(ctx context.Context) (string, error) {
	value, ok := metainfo.GetValue(ctx, "token")
	if !ok {
		return "", nil
	}
	return value, nil
}

func GetMetaInfoCurrentUserId(ctx context.Context) (int64, error) {
	value, ok := metainfo.GetPersistentValue(ctx, constant.CurrentUserId)
	if !ok {
		return 0, nil
	}
	currentUserId, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return 0, handler.HandleError(err)
	}
	return currentUserId, nil
}

func GetMetaInfoCurrentUserName(ctx context.Context) (string, error) {
	value, ok := metainfo.GetPersistentValue(ctx, constant.CurrentUserName)
	if !ok {
		return "", nil
	}
	return value, nil
}

func GetMetaInfoCurrentOrganizationId(ctx context.Context) (int64, error) {
	value, ok := metainfo.GetPersistentValue(ctx, constant.CurrentOrganizationId)
	if !ok {
		return 0, nil
	}
	currentOrganizationId, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return 0, handler.HandleError(err)
	}
	return currentOrganizationId, nil
}

func GetMetaInfoAuthElements(ctx context.Context) ([]string, error) {
	value, ok := metainfo.GetPersistentValue(ctx, constant.AuthElements)
	if !ok {
		return nil, nil
	}
	var result []string
	err := json.Unmarshal([]byte(value), &result)
	return result, handler.HandleError(err)
}
