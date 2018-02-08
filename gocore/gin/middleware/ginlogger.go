package middleware

import (
	"time"

	"github.bus.zalan.do/ale/gocore/logger"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		var flowid string

		if fid, ok := c.Keys[FlowIDKey]; ok {
			if fid2, ok := fid.(string); ok {
				flowid = fid2
			}
		}

		logger.LogInfo(flowid, "%d ; call to %s %s took %13v ; IP was : %s ; comment : %s",
			statusCode,
			method,
			path,
			latency,
			clientIP,
			comment,
		)

	}
}
