package middlewares

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	pctx "github.com/luoye-g/webgate/pkg/ctx"
	"github.com/luoye-g/webgate/pkg/log"
)

const traceIDHeader = "X-Trace-Id"

// bodyLogWriter captures response body for logging (optional usage).
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func genTraceID() string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return time.Now().Format("20060102150405.000000")
	}
	return hex.EncodeToString(b[:])
}

// Logger returns a middleware that logs each HTTP request with slog.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// trace id: from header or generate
		traceID := c.GetHeader(traceIDHeader)
		if traceID == "" {
			traceID = genTraceID()
		}
		pctx.SetTraceID(c, traceID)
		c.Writer.Header().Set(traceIDHeader, traceID)

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		attrs := []any{
			slog.String("trace_id", traceID),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("client_ip", c.ClientIP()),
			slog.Int("status", status),
			slog.Float64("latency_ms", float64(latency.Microseconds())/1000.0),
			slog.Int("body_size", c.Writer.Size()),
		}
		if uid := pctx.GetUserID(c); uid != 0 {
			attrs = append(attrs, slog.Uint64("user_id", uid))
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.String()))
		}

		msg := "http_request"
		switch {
		case status >= http.StatusInternalServerError:
			log.L().Error(msg, attrs...)
		case status >= http.StatusBadRequest:
			log.L().Warn(msg, attrs...)
		default:
			log.L().Info(msg, attrs...)
		}
	}
}

// Recovery returns a middleware that recovers from panics and logs them.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				traceID := pctx.GetTraceID(c)
				log.L().Error("panic_recovered",
					slog.String("trace_id", traceID),
					slog.String("path", c.Request.URL.Path),
					slog.Any("panic", r),
					slog.String("stack", string(debug.Stack())),
				)
				if !c.Writer.Written() {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"code": 500,
						"msg":  "internal error",
					})
				} else {
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}

// ensure io import used (kept for potential future body capture usage)
var _ = io.Discard
