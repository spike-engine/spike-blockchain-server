package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerToFile() gin.HandlerFunc {
	fileName := os.Getenv("LOG_FILE_PREFIX") + " " + time.Now().Format("2021-01-01 12.01.01")

	src, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	logger := logrus.New()
	logger.Out = src
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		ip := c.ClientIP()
		logger.WithFields(logrus.Fields{
			"code":     statusCode,
			"method":   reqMethod,
			"uri":      reqUri,
			"duration": duration.Seconds(),
			"ip":       ip,
			"headers":  c.Request.Header,
		}).Info(c.Writer.Header().Get("Content-Type"))
	}
}
