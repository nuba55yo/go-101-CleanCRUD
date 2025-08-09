package logging

import (
	"context"

	"github.com/nuba55yo/go-101-CleanCRUD/application/interfaces"
	"go.uber.org/zap"
)

// ZapLogger คืออแดปเตอร์ logger ที่ implement interfaces.Logger
// ใช้ zap แบบ production (ปรับ config เพิ่มเติมได้ตามต้องการ)
type ZapLogger struct {
	inner *zap.SugaredLogger
}

func NewZapLogger() (interfaces.Logger, func() error, error) {
	z, err := zap.NewProduction()
	if err != nil {
		return nil, nil, err
	}
	return &ZapLogger{inner: z.Sugar()}, z.Sync, nil
}

func (logger *ZapLogger) Info(requestContext context.Context, message string, keyValues ...any) {
	logger.inner.Infow(message, keyValues...)
}

func (logger *ZapLogger) Warn(requestContext context.Context, message string, keyValues ...any) {
	logger.inner.Warnw(message, keyValues...)
}

func (logger *ZapLogger) Error(requestContext context.Context, message string, keyValues ...any) {
	logger.inner.Errorw(message, keyValues...)
}
