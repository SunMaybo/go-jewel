package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
)

var LOGGER *zap.Logger

func GetLog(level string) {
	var atomicLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	switch strings.ToLower(level) {
	case "info":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		atomicLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	cfg := zap.Config{
		Encoding:         "console",
		Level:            atomicLevel,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
	LOGGER = logger
}
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		switch {
		case statusCode >= 400 && statusCode <= 499:
			{
				logger.Warn("[GIN]",
					zap.Int("statusCode", statusCode),
					zap.String("latency", latency.String()),
					zap.String("clientIP", clientIP),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("error", c.Errors.String()),
				)
			}
		case statusCode >= 500:
			{
				logger.Error("[GIN]",
					//zap.String("statusColor", statusColor),
					zap.Int("statusCode", statusCode),
					zap.String("latency", latency.String()),
					zap.String("clientIP", clientIP),
					//zap.String("methodColor", methodColor),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("error", c.Errors.String()),
				)
			}
		default:
			logger.Info("[GIN]",
				//zap.String("statusColor", statusColor),
				zap.Int("statusCode", statusCode),
				zap.String("latency", latency.String()),
				zap.String("clientIP", clientIP),
				//zap.String("methodColor", methodColor),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("error", c.Errors.String()),
			)
		}
	}
}
