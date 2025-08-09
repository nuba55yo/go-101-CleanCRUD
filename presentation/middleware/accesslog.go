package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	logMutex      sync.Mutex
	currentBucket string
	currentFile   *os.File
)

type responseRecorder struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *responseRecorder) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w *responseRecorder) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func timeBucket(t time.Time) string {
	b := (t.Minute() / 10) * 10 // 10 นาที/ไฟล์
	return fmt.Sprintf("%04d-%02d-%02d_%02d-%02d", t.Year(), t.Month(), t.Day(), t.Hour(), b)
}

func ensureLogFile(now time.Time) (*os.File, error) {
	logMutex.Lock()
	defer logMutex.Unlock()

	bucket := timeBucket(now)
	if currentFile != nil && bucket == currentBucket {
		return currentFile, nil
	}
	if currentFile != nil {
		_ = currentFile.Close()
	}

	dateFolder := now.Format("2006-01-02")
	dir := filepath.Join("logs", dateFolder)
	_ = os.MkdirAll(dir, 0o755)

	name := fmt.Sprintf("log_%s.log", bucket)
	full := filepath.Join(dir, name)

	f, err := os.OpenFile(full, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	currentBucket = bucket
	currentFile = f
	return f, nil
}

func writeLine(now time.Time, module, level, message string) {
	f, err := ensureLogFile(now)
	if err != nil {
		return
	}
	line := fmt.Sprintf("%s [%s] [%s] %s\n", now.Format("2006-01-02 15:04:05.000"), module, level, message)
	_, _ = f.WriteString(line)
}

func readRequestBodySafely(r *http.Request, limit int64) string {
	if r.Body == nil {
		return ""
	}
	body, _ := io.ReadAll(io.LimitReader(r.Body, limit))
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return string(body)
}

func moduleFromRoute(fullPath string) string {
	if fullPath == "" {
		return "-"
	}
	path := strings.TrimPrefix(fullPath, "/")
	parts := strings.Split(path, "/")
	// ex: /api/v1/books → books
	if len(parts) >= 3 && parts[0] == "api" && strings.HasPrefix(parts[1], "v") {
		return parts[2]
	}
	if len(parts) > 0 {
		return parts[0]
	}
	return "-"
}

func levelFromStatus(status int) string {
	switch {
	case status >= 500:
		return "error"
	case status >= 400:
		return "warn"
	default:
		return "info"
	}
}

// AccessLog: บันทึก request/response + หมุนไฟล์ทุก 10 นาที
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqBody := readRequestBodySafely(c.Request, 1<<20) // 1MB

		rec := &responseRecorder{ResponseWriter: c.Writer}
		c.Writer = rec

		c.Next()

		latency := time.Since(start)
		status := rec.Status()
		route := c.FullPath()
		module := moduleFromRoute(route)

		msg := fmt.Sprintf(
			"status=%d method=%s route=%s ip=%s latency=%s req=%s res=%s",
			status,
			c.Request.Method,
			route,
			c.ClientIP(),
			latency.String(),
			strings.ReplaceAll(reqBody, "\n", " "),
			strings.ReplaceAll(rec.body.String(), "\n", " "),
		)
		writeLine(time.Now(), module, levelFromStatus(status), msg)
	}
}
