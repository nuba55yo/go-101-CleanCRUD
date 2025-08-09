package interfaces

import "context"

// Logger คือสัญญา logging ที่ไม่ผูกกับ framework
// เลเยอร์ infrastructure จะทำตัวจริง (เช่น zap) มา implement อันนี้
type Logger interface {
	Info(requestContext context.Context, message string, keyValues ...any)
	Warn(requestContext context.Context, message string, keyValues ...any)
	Error(requestContext context.Context, message string, keyValues ...any)
}
