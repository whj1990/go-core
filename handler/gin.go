package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type BasePaginationRequest struct {
	PageNum        int32  `json:"pageNum" form:"pageNum"`
	PageSize       int32  `json:"pageSize" form:"pageSize"`
	SortField      string `json:"sortField" form:"sortField"`
	SortOrder      int32  `json:"sortOrder" form:"sortOrder"`
	OrganizationId int64  `json:"organizationId" form:"organizationId"`
}

type BaseRequest struct {
	Id             int64 `json:"id" form:"id"`
	OrganizationId int64 `json:"organizationId" form:"organizationId"`
}

type BaseResponse struct {
	Id             int64     `json:"id"`
	OrganizationId int64     `json:"organizationId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
	Code  int         `json:"code"`
}

type listResponse struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
	Code  int         `json:"code"`
}

func HandleUnauthorizedResponse(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, response{
		Error: err.Error(),
	})
}

func HandleErrorResponse(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, response{
		Error: err.Error(),
	})
}

func HandleSuccessResponse(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, response{
		Data: "ok",
	})
}

func HandleSuccessDataResponse(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, response{
		Data: data,
	})
}

func HandleSuccessListDataResponse(c *gin.Context, data interface{}, total int64) {
	c.AbortWithStatusJSON(http.StatusOK, listResponse{
		Data:  data,
		Total: total,
	})
}
