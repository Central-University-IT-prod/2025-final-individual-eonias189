package logger

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

var once sync.Once

var logger *zap.Logger

func Get(level string) (*zap.Logger, error) {
	zapLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)
		logLevel := zap.NewAtomicLevelAt(zapLevel)

		cfg := zap.NewProductionEncoderConfig()
		cfg.TimeKey = "timestamp"
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(cfg)

		core := zapcore.NewCore(consoleEncoder, stdout, logLevel)

		logger = zap.New(
			core,
			zap.AddCaller(),
		)
	})

	return logger, nil
}

func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return zap.NewNop()
}

func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == l {
			// Do not store same logger.
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}
