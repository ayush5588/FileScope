package internal

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetLogger ...
func GetLogger() *zap.SugaredLogger {
	// Setup logger
	zapProdConfig := zap.NewProductionConfig()
	// Modify the logger to show rfc3339 date & time format
	zapProdConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	zapProd, _ := zapProdConfig.Build()

	logger := zapProd.Sugar()

	return logger
}
