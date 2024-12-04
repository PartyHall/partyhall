package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		if query == "/.well-known/mercure" {
			return
		}

		method := c.Request.Method
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		latency := time.Since(start)

		msg := "Query: " + query

		var level zapcore.Level
		switch {
		case status >= 500:
			level = zapcore.ErrorLevel
		case status >= 400:
			level = zapcore.WarnLevel
		default:
			level = zapcore.InfoLevel
		}

		if len(c.Errors) > 0 {
			errorMsgs := make([]string, len(c.Errors))
			for i, err := range c.Errors {
				errorMsgs[i] = err.Error()
			}
			msg = fmt.Sprintf("%s\nErrors:\n%s",
				msg,
				strings.Join(errorMsgs, "\n"),
			)
		}

		Message(level, msg, []any{
			"method", method,
			"status", status,
			"latency", latency,
			"ip", clientIP,
			"path", path,
			"query", query,
		})
	}
}
