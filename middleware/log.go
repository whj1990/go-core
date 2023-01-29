package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

func ApiLogMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		loggerMap := make(map[string]interface{})
		loggerMap["clientIP"] = param.ClientIP
		loggerMap["time"] = param.TimeStamp.Format(time.RFC1123)
		loggerMap["method"] = param.Method
		loggerMap["path"] = param.Path
		loggerMap["proto"] = param.Request.Proto
		loggerMap["statusCode"] = param.StatusCode
		loggerMap["latency"] = param.Latency
		loggerMap["userAgent"] = param.Request.UserAgent()
		loggerMap["errorMessage"] = param.ErrorMessage
		loggerJson, _ := json.Marshal(loggerMap)
		return string(loggerJson)
	})
}
