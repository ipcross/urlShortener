package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Initialize(level string) (*zap.Logger, error) {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("zap.ParseAtomicLeve: %w", err)
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("cfg.Build: %w", err)
	}

	return zl, nil
}

func Sync(l *zap.Logger) {
	_ = l.Sync()
}

func RequestLogger(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				l.Info("RequestLogger",
					zap.String("proto", r.Proto),
					zap.String("url", r.URL.Path),
					zap.String("method", r.Method),
					zap.Duration("duration", time.Since(t1)),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()))
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
