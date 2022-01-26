package modlog

import (
	"context"
	"net/http"

	zl "github.com/rs/zerolog"
)

type Logger struct {
	module string
	logger zl.Logger
}

func (l *Logger) SetLevel(level zl.Level) {
	l.logger = l.logger.Level(level)
}

func (l *Logger) Event(eventName string) *zl.Event {
	return l.logger.Info().Str(keyEventType, eventName)
}

func (l *Logger) EventContext(ctx context.Context, eventName string) *zl.Event {
	return withContext(l.Event(eventName), ctx)
}

func (l *Logger) LogError(err error, msg string) {
	if err == nil {
		return
	}
	l.Error().Err(err).Msg(msg)
}

func (l *Logger) Error() *zl.Event {
	return l.logger.Error()
}

func (l *Logger) ErrorRequest(r *http.Request) *zl.Event {
	return withContext(l.logger.Error(), r.Context())
}

func (l *Logger) ErrorContext(ctx context.Context) *zl.Event {
	return withContext(l.logger.Error(), ctx)
}

func (l *Logger) Warn() *zl.Event {
	return l.logger.Warn()
}

func (l *Logger) WarnRequest(r *http.Request) *zl.Event {
	return withContext(l.logger.Warn(), r.Context())
}

func (l *Logger) WarnContext(ctx context.Context) *zl.Event {
	return withContext(l.logger.Warn(), ctx)
}

func (l *Logger) Info() *zl.Event {
	return l.logger.Info()
}

func (l *Logger) InfoRequest(r *http.Request) *zl.Event {
	return withContext(l.logger.Info(), r.Context())
}

func (l *Logger) InfoContext(ctx context.Context) *zl.Event {
	return withContext(l.logger.Info(), ctx)
}

func (l *Logger) Debug() *zl.Event {
	return l.logger.Debug()
}

func (l *Logger) DebugRequest(r *http.Request) *zl.Event {
	return withContext(l.logger.Debug(), r.Context())
}

func (l *Logger) DebugContext(ctx context.Context) *zl.Event {
	return withContext(l.logger.Debug(), ctx)
}

func (l *Logger) Trace() *zl.Event {
	return l.logger.Trace()
}

func (l *Logger) TraceRequest(r *http.Request) *zl.Event {
	return withContext(l.logger.Trace(), r.Context())
}

func (l *Logger) TraceContext(ctx context.Context) *zl.Event {
	return withContext(l.logger.Trace(), ctx)
}

func withContext(evt *zl.Event, ctx context.Context) *zl.Event {
	for _, fn := range global.contextHooks {
		evt = fn(evt, ctx)
	}
	return evt
}
