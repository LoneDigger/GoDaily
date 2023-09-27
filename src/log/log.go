package log

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 存放紀錄上限
const limit = 1000

var LogHistory *Logger

type Logger struct {
	L     *logrus.Logger
	Sb    strings.Builder
	Array []string
	Func  gin.HandlerFunc
}

func init() {
	LogHistory = &Logger{}
	LogHistory.L = logrus.New()
	LogHistory.L.Out = LogHistory

	LogHistory.L.SetLevel(logrus.InfoLevel)

	LogHistory.L.SetFormatter(&logrus.TextFormatter{
		DisableSorting:   false,
		DisableTimestamp: false,
		TimestampFormat:  "2006-01-02 15:04:05",
	})

	LogHistory.Func = func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		path := c.Request.URL.Path

		switch {
		case !strings.HasPrefix(path, "/public"):
			endTime := time.Now()
			latencyTime := endTime.Sub(startTime)
			reqMethod := c.Request.Method
			statusCode := c.Writer.Status()
			clientIP := c.ClientIP()
			bodySize := c.Writer.Size()
			code := c.GetString("code")

			raw := c.Request.URL.RawQuery
			if raw != "" {
				path = path + "?" + raw
			}

			LogHistory.L.WithFields(logrus.Fields{
				"Method":   reqMethod,
				"Ip":       clientIP,
				"Status":   statusCode,
				"Code":     code,
				"BodySize": bodySize,
				"Latency":  latencyTime,
				"Path":     path,
			}).Info("Gin")
		}
	}
}

func (l *Logger) String() string {
	return strings.Join(l.Array, "")
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.Sb.Write(p)

	// 有上限存放紀錄
	if len(l.Array) > limit {
		// remove last
		l.Array = l.Array[:len(l.Array)-1]
	}

	// add first
	l.Array = append([]string{l.Sb.String()}, l.Array...)
	l.Sb.Reset()

	return os.Stdout.Write(p)
}
